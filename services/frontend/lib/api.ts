const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '/api'

// Types
export interface PerfumeCharacteristics {
    density: number
    earthiness: number
    floralcy: number
    freshness: number
    fruityness: number
    powderiness: number
    spiciness: number
    sweetness: number
    warmth: number
    woodiness: number
}

export interface PerfumeProperties {
    perfume_type: string
    family: string[]
    upper_notes: string[]
    core_notes: string[]
    base_notes: string[]
    tags?: string[]
    upper_characteristics?: PerfumeCharacteristics
    core_characteristics?: PerfumeCharacteristics
    base_characteristics?: PerfumeCharacteristics
}

export interface PerfumeVariant {
    volume: number
    link: string
    price: number
}

export interface ShopInfo {
    shop_name: string
    domain: string
    variants: PerfumeVariant[]
}

export interface Perfume {
    id?: number
    brand: string
    name: string
    sex: string
    image_url: string
    properties: PerfumeProperties
    shops: ShopInfo[]

    rank?: number
    similarity_score?: number
}

export interface PerfumeResponse {
    perfumes: Perfume[]
    total: number
}

export interface SuggestionRequest {
    brand: string
    name: string
    use_ai?: boolean
    sex?: 'male' | 'unisex' | 'female'
}

export interface SuggestionResponse {
    suggested: Array<{
        perfume: {
            brand: string
            name: string
            sex: string
            image_url: string
            properties: PerfumeProperties
            shops: ShopInfo[]
        }
        rank: number
        similarity_score: number
    }>
}

interface ErrorResponse {
    error: string
    message: string
}

// API Client class
class APIClient {
    private async request<T>(
        url: string,
        options: RequestInit = {}
    ): Promise<T> {
        const response = await fetch(url, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            ...options,
        })

        if (response.status === 204) {
            return { suggested: [] } as T
        }

        if (!response.ok) {
            let errorMessage = `API request failed: ${response.status} ${response.statusText}`
            try {
                const errorData: ErrorResponse = await response.json()
                errorMessage = errorData.message || errorData.error || errorMessage
            } catch {
            }
            throw new Error(errorMessage)
        }

        return response.json()
    }

    async getSuggestions(request: SuggestionRequest): Promise<SuggestionResponse> {
        const params = new URLSearchParams({
            brand: request.brand,
            name: request.name,
        })

        if (request.use_ai !== undefined) {
            params.append('use_ai', request.use_ai.toString())
        }

        if (request.sex !== undefined) {
            params.append('sex', request.sex)
        }

        const url = `${API_BASE_URL}/perfume/suggest?${params.toString()}`
        return this.request<SuggestionResponse>(url)
    }
}

// Export singleton instance
export const apiClient = new APIClient()

// Export individual functions for convenience
export const getSuggestions = (request: SuggestionRequest) =>
    apiClient.getSuggestions(request)
