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
                description: `${item.perfume.properties.type} • ${item.perfume.properties.sex}`,
                notes: {
                    upper: item.perfume.properties.upper_notes,
                    middle: item.perfume.properties.middle_notes,
                    base: item.perfume.properties.base_notes
                },
                type: item.perfume.properties.type,
                family: item.perfume.properties.family,
                sex: item.perfume.properties.sex,
                links: item.perfume.links,
                image: item.perfume.image_url
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
