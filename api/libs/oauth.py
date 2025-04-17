import json
import logging  # 二开部分，针对oa登录报错问题，记录返回的code
import urllib.parse
from dataclasses import dataclass
from typing import Optional

import requests
from configs import dify_config  # Extend OAuto third-party login
from extensions.ext_database import db  # Extend OAuto third-party login
from models.system_extend import SystemIntegrationExtend, SystemIntegrationClassify  # Extend OAuto third-party login


@dataclass
class OAuthUserInfo:
    id: str
    name: str
    email: str


class OAuth:
    def __init__(self, client_id: str, client_secret: str, redirect_uri: str):
        self.client_id = client_id
        self.client_secret = client_secret
        self.redirect_uri = redirect_uri

    def get_authorization_url(self):
        raise NotImplementedError()

    def get_access_token(self, code: str):
        raise NotImplementedError()

    def get_raw_user_info(self, token: str):
        raise NotImplementedError()

    def get_user_info(self, token: str) -> OAuthUserInfo:
        raw_info = self.get_raw_user_info(token)
        return self._transform_user_info(raw_info)

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        raise NotImplementedError()


class GitHubOAuth(OAuth):
    _AUTH_URL = "https://github.com/login/oauth/authorize"
    _TOKEN_URL = "https://github.com/login/oauth/access_token"
    _USER_INFO_URL = "https://api.github.com/user"
    _EMAIL_INFO_URL = "https://api.github.com/user/emails"

    def get_authorization_url(self, invite_token: Optional[str] = None):
        params = {
            "client_id": self.client_id,
            "redirect_uri": self.redirect_uri,
            "scope": "user:email",  # Request only basic user information
        }
        if invite_token:
            params["state"] = invite_token
        return f"{self._AUTH_URL}?{urllib.parse.urlencode(params)}"

    def get_access_token(self, code: str):
        data = {
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            "code": code,
            "redirect_uri": self.redirect_uri,
        }
        headers = {"Accept": "application/json"}
        response = requests.post(self._TOKEN_URL, data=data, headers=headers)

        response_json = response.json()
        access_token = response_json.get("access_token")

        if not access_token:
            raise ValueError(f"Error in GitHub OAuth: {response_json}")

        return access_token

    def get_raw_user_info(self, token: str):
        headers = {"Authorization": f"token {token}"}
        response = requests.get(self._USER_INFO_URL, headers=headers)
        response.raise_for_status()
        user_info = response.json()

        email_response = requests.get(self._EMAIL_INFO_URL, headers=headers)
        email_info = email_response.json()
        primary_email: dict = next((email for email in email_info if email["primary"] == True), {})

        return {**user_info, "email": primary_email.get("email", "")}

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        email = raw_info.get("email")
        if not email:
            email = f"{raw_info['id']}+{raw_info['login']}@users.noreply.github.com"
        return OAuthUserInfo(id=str(raw_info["id"]), name=raw_info["name"], email=email)


class GoogleOAuth(OAuth):
    _AUTH_URL = "https://accounts.google.com/o/oauth2/v2/auth"
    _TOKEN_URL = "https://oauth2.googleapis.com/token"
    _USER_INFO_URL = "https://www.googleapis.com/oauth2/v3/userinfo"

    def get_authorization_url(self, invite_token: Optional[str] = None):
        params = {
            "client_id": self.client_id,
            "response_type": "code",
            "redirect_uri": self.redirect_uri,
            "scope": "openid email",
        }
        if invite_token:
            params["state"] = invite_token
        return f"{self._AUTH_URL}?{urllib.parse.urlencode(params)}"

    def get_access_token(self, code: str):
        data = {
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            "code": code,
            "grant_type": "authorization_code",
            "redirect_uri": self.redirect_uri,
        }
        headers = {"Accept": "application/json"}
        response = requests.post(self._TOKEN_URL, data=data, headers=headers)

        response_json = response.json()
        access_token = response_json.get("access_token")

        if not access_token:
            raise ValueError(f"Error in Google OAuth: {response_json}")

        return access_token

    def get_raw_user_info(self, token: str):
        headers = {"Authorization": f"Bearer {token}"}
        response = requests.get(self._USER_INFO_URL, headers=headers)
        response.raise_for_status()
        return response.json()

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        return OAuthUserInfo(id=str(raw_info["sub"]), name="", email=raw_info["email"])


# Extend Start: OAuth2
class OaOAuth(OAuth):

    def get_auto2_conf(self):
        # oauth start
        integration = db.session.query(SystemIntegrationExtend).filter(
            SystemIntegrationExtend.classify == SystemIntegrationClassify.SYSTEM_INTEGRATION_OAUTH_TWO).first()
        if integration is None or (integration and not integration.status):
            return {
                "integration": integration,
                "config": {},
                "passwd": ""
            }
        return {
            "integration": integration,
            "passwd": integration.decodeSecret(),
            "config": json.loads(integration.config)
        }

    def extract_data(self, dictionary, path):
        """
        从字典中提取指定路径的数据
        支持通配符'*'获取列表中所有元素的特定字段

        Args:
            dictionary (dict): 源字典
            path (str): 以点分隔的路径，如 "data.info.name" 或 "data.items.*.name"

        Returns:
            提取的数据
        """
        parts = path.split('.')
        current = dictionary

        for i, part in enumerate(parts):
            if part == '*' and isinstance(current, list):
                # 处理列表中的每个元素
                remainder = '.'.join(parts[i + 1:])
                if remainder:
                    return [self.extract_data(item, remainder) for item in current]
                else:
                    return current
            elif isinstance(current, dict) and part in current:
                current = current[part]
            else:
                return None

        return current

    def get_authorization_url(self, invite_token: Optional[str] = None):
        auto2_conf = self.get_auto2_conf()
        integration = auto2_conf.get('integration')
        if integration is None:
            return
        # 构建查询参数
        query_string = urllib.parse.urlencode({
            'redirect_uri': dify_config.CONSOLE_API_URL + "/console/api/oauth/authorize/oauth2",
            'client_id': integration.app_id,
        })

        # 构建完整URL
        config = auto2_conf.get('config')
        return f"{config.get('server_url')}{config.get('authorize_url')}?{query_string}"

    def get_access_token(self, code: str):
        auto2_conf = self.get_auto2_conf()
        integration = auto2_conf.get('integration')
        if integration is None:
            return ""
        data = {
            "code": code,
            "client_id": integration.app_id,
            "grant_type": "authorization_code",
            "client_secret": auto2_conf.get('passwd'),
            "redirect_uri": dify_config.CONSOLE_API_URL + "/console/api/oauth/authorize/oauth2",
        }
        headers = {"Accept": "application/json"}
        config = auto2_conf.get('config')
        response = requests.post(f"{config.get('server_url')}{config.get('token_url')}", data=data, headers=headers)
        response.encoding = "utf-8"
        if response.status_code != 200:
            return ""
        try:
            response_json = response.json()
        except:
            return ""
        access_token = response_json.get("access_token")

        return access_token

    def get_raw_user_info(self, token: str):
        auto2_conf = self.get_auto2_conf()
        if auto2_conf.get('integration') is None:
            return ""
        config = auto2_conf.get('config')
        headers = {"Authorization": f"Bearer {token}"}
        response = requests.get(f"{config.get('server_url')}{config.get('userinfo_url')}", headers=headers)
        response.raise_for_status()
        return response.json()

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        # 检查 raw_info 是否为空或为 None
        auto2_conf = self.get_auto2_conf()
        if not raw_info or not isinstance(raw_info, dict) or auto2_conf.get('integration') is None:
            return OAuthUserInfo(
                id="",
                name="",
                email="",
            )
        # 提取参数
        config = auto2_conf.get('config')
        name = self.extract_data(raw_info, config.get('user_name_field'))
        email = self.extract_data(raw_info, config.get('user_email_field'))
        username = self.extract_data(raw_info, config.get('user_id_field'))
        if not username:
            raise ValueError("OAuth2返回用户数据格式不正确。请返回进行重新登录。")

        return OAuthUserInfo(
            id=str(username) if username is not None else None,
            name=str(name) if name is not None else None,
            email=str(email) if email is not None else None,
        )
# Extend Stop: OAuth
