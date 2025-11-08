import { useState, useCallback } from 'react'
import { apiClient, Perfume, SuggestionRequest } from '@/lib/api'

interface UsePerfumeAPIState {
    loading: boolean
    error: string | null
    data: Perfume[] | null
}

export const usePerfumeAPI = () => {
    const [state, setState] = useState<UsePerfumeAPIState>({
        loading: false,
        error: null,
        data: null,
    })

    const reset = useCallback(() => {
        setState({
            loading: false,
            error: null,
            data: null,
        })
    }, [])

    const getSuggestions = useCallback(async (request: SuggestionRequest) => {
        setState(prev => ({ ...prev, loading: true, error: null }))

        try {
            const response = await apiClient.getSuggestions(request)
            const perfumes: Perfume[] = response.suggested.map(item => ({
                id: item.rank,
                brand: item.perfume.brand,
                name: item.perfume.name,
                sex: item.perfume.sex,
                image_url: item.perfume.image_url,
                properties: item.perfume.properties,
                shops: item.perfume.shops,
                rank: item.rank,
                similarity_score: item.similarity_score
            }))

            setState({
                loading: false,
                error: null,
                data: perfumes,
            })
            return perfumes
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'Failed to get suggestions'
            setState({
                loading: false,
                error: errorMessage,
                data: null,
            })
            throw error
        }
    }, [])

    return {
        ...state,
        getSuggestions,
        reset,
    }
}
