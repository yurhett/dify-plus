'use client'
import React, { useEffect, useState } from 'react'
import { fetchUserMoney } from '@/service/common-extend'
import type { UserMoney } from '@/models/common-extend'
import cn from 'classnames'

const AccountMoneyExtend = () => {
  const [userMoney, setUserMoney] = useState<UserMoney>({ used_quota: 0, total_quota: 15 })
  const [isFetched, setIsFetched] = useState(false)
  const exchangeRate = 7.26 // 美元转人民币固定汇率

  const getUserMoney = async () => {
    const data: any = await fetchUserMoney()
    setUserMoney(data)
  }

  useEffect(() => {
    getUserMoney()
    setIsFetched(true)
  }, [])

  if (!isFetched)
    return null

  // 计算额度（确保使用数字类型）
  const usedQuota = Number(userMoney.used_quota) || 0
  const totalQuota = Number(userMoney.total_quota) || 15
  const remainingQuota = totalQuota - usedQuota

  // 转换为人民币并保留2位小数
  const usedRMB = (usedQuota * exchangeRate).toFixed(2)
  const totalRMB = (totalQuota * exchangeRate).toFixed(2)
  const remainingRMB = (remainingQuota * exchangeRate).toFixed(2)

  // 判断警示级别
  const isRedAlert = Number(remainingRMB) < 10 // 余额不足10元人民币，显示红色
  const isYellowAlert = Number(usedRMB) > 50 && !isRedAlert // 使用超过50元人民币，显示黄色

  // 根据警示级别设置颜色
  const alertColorClass = isRedAlert
    ? 'text-red-500'
    : isYellowAlert
      ? 'text-yellow-500'
      : 'text-gray-700'

  return (
    <div
      rel='noopener noreferrer'
      className='flex items-center overflow-hidden rounded-md border border-gray-200 text-xs leading-[18px]'
    >
      <div className='flex items-center bg-gray-100 px-2 py-1 font-medium'>
        额度
      </div>
      <div className='flex items-center border-l border-gray-200 bg-white px-2 py-1.5'>
        <span className='mr-1 text-gray-600'>已用:</span>
        <span
          className={cn(
            'font-bold transition-all duration-300',
            alertColorClass,
            'text-sm md:text-base', // 默认字体稍大，响应式设计
          )}
        >
          ¥{usedRMB}
        </span>
        <span className='mx-1 text-gray-400'>/</span>
        <span className='text-gray-500'>
          ¥{totalRMB.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
        </span>
      </div>
    </div>
  )
}

export default AccountMoneyExtend
