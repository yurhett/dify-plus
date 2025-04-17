export enum SSOProtocol {
  SAML = 'saml',
  OIDC = 'oidc',
  OAuth2 = 'oauth2',
}

export enum LicenseStatus {
  NONE = 'none',
  INACTIVE = 'inactive',
  ACTIVE = 'active',
  EXPIRING = 'expiring',
  EXPIRED = 'expired',
  LOST = 'lost',
}

type License = {
  status: LicenseStatus
  expired_at: string | null
}

export type SystemFeatures = {
  sso_enforced_for_signin: boolean
  sso_enforced_for_signin_protocol: SSOProtocol | ''
  sso_enforced_for_web: boolean
  sso_enforced_for_web_protocol: SSOProtocol | ''
  enable_web_sso_switch_component: boolean
  enable_marketplace: boolean
  enable_email_code_login: boolean
  enable_email_password_login: boolean
  enable_social_oauth_login: boolean
  is_allow_create_workspace: boolean
  is_allow_register: boolean
  is_email_setup: boolean
  license: License
  is_custom_auth2: string // extend: Customizing AUTH2
  is_custom_auth2_logout: string // extend: AUTH2 logout url
  ding_talk_client_id: string // Extend: DingTalk third-party login
  ding_talk_corp_id: string // Extend: DingTalk sidebar login
  ding_talk: boolean // Extend: switch DingTalk sidebar login
}

export const defaultSystemFeatures: SystemFeatures = {
  sso_enforced_for_signin: false,
  sso_enforced_for_signin_protocol: '',
  sso_enforced_for_web: false,
  sso_enforced_for_web_protocol: '',
  enable_web_sso_switch_component: false,
  enable_marketplace: false,
  enable_email_code_login: false,
  enable_email_password_login: false,
  enable_social_oauth_login: false,
  is_allow_create_workspace: false,
  is_allow_register: false,
  is_email_setup: false,
  license: {
    status: LicenseStatus.NONE,
    expired_at: '',
  },
  is_custom_auth2: '', // extend: Customizing AUTH2
  ding_talk_client_id: '', // Extend: DingTalk third-party login
  ding_talk_corp_id: '', // Extend: DingTalk sidebar login
  ding_talk: false, // Extend: switch DingTalk sidebar login
}
