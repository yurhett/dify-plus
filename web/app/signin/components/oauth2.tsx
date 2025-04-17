import React from 'react'
import { useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import style from '../page.module.css'
import Button from '@/app/components/base/button'
import classNames from '@/utils/classnames'
import { apiPrefix } from '@/config'

export default function OAuth2() {
  const { t } = useTranslation()
  const router = useRouter()

  /* Extend: start 钉钉快捷登录按钮 */
  const OAuth2Login = () => {
    router.replace(`${apiPrefix}/oauth/login/oauth2`)
  }

  return <>
    <div className="mb-2">
      <a onClick={OAuth2Login}>
        <Button
          className="w-full"
        >
          <span className={
            classNames(
              style.oauth2Icon,
              'w-5 h-5 mr-2',
            )
          }
          />
          <span className="truncate">{t('appOverview.overview.appInfo.settings.sso.label')}</span>
        </Button>
      </a>
    </div>
  </>
}
