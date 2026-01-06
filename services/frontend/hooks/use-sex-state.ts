import { useState, useEffect } from 'react'

const SEX_STORAGE_KEY = 'scently_selected_sex'

type SexValue = 'male' | 'unisex' | 'female'

export function useSexState(defaultValue: SexValue = 'unisex') {
    const [sex, setSex] = useState<SexValue>(() => {
        if (typeof window !== 'undefined') {
            const stored = localStorage.getItem(SEX_STORAGE_KEY)
            if (stored && ['male', 'unisex', 'female'].includes(stored)) {
                return stored as SexValue
            }
        }
        return defaultValue
    })

    useEffect(() => {
        localStorage.setItem(SEX_STORAGE_KEY, sex)
    }, [sex])

    return [sex, setSex] as const
}

