import React from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import Button from '@/app/components/base/button'
import { fetchWebOAuth2SSOUrl } from '@/service/share'

const OAuth2 = () => {
  const { t } = useTranslation()
  const router = useRouter()
  const searchParams = useSearchParams()

  const redirectUrl = searchParams.get('redirect_url')

  const getAppCodeFromRedirectUrl = () => {
    if (!redirectUrl)
      return null
    const url = new URL(`${window.location.origin}${decodeURIComponent(redirectUrl)}`)
    const appCode = url.pathname.split('/').pop()
    if (!appCode)
      return null

    return appCode
  }

  const OAuth2Login = async () => {
    const appCode = getAppCodeFromRedirectUrl()
    if (!redirectUrl || !appCode) {
      console.error('Invalid redirect URL or app code')
      return
    }

    try {
      const response = await fetchWebOAuth2SSOUrl(appCode, redirectUrl)
      router.push(response.url)
    }
    catch (error) {
      console.error('OAuth2 login error:', error)
    }
  }

  return (
    <div className="mb-2">
      <Button
        className="w-full"
        onClick={OAuth2Login}
      >
        <span className="truncate">{t('appOverview.overview.appInfo.settings.sso.label')}</span>
      </Button>
    </div>
  )
}

export default OAuth2
