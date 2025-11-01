// API configuration - use Next.js proxy to avoid CORS issues
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '/api'

// Types
export interface Perfume {
    id: number
    brand: string
    name: string
    description?: string
    image?: string
    notes?: {
        upper: string[]
        middle: string[]
        base: string[]
    }
    type?: string
    family?: string[]
    sex?: string
    links?: Record<string, string>
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
    input: {
        brand: string
        name: string
        ok: boolean
        advise_type: string
        sex: string
    }
    suggested: Array<{
        perfume: {
            brand: string
            name: string
            sex: string
            properties: {
                type: string
                family: string[]
                upper_notes: string[]
                middle_notes: string[]
                base_notes: string[]
            }
            links: Record<string, string>
            image_url: string
        }
        rank: number
        similarity_score: number
    }>
    success: boolean
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

        if (!response.ok) {
            throw new Error(`API request failed: ${response.status} ${response.statusText}`)
        }

        return response.json()
    }

    // Get perfume suggestions from perfumist service
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

        const url = `${API_BASE_URL}/v1/suggest/perfume?${params.toString()}`
        const origin =
            typeof window !== "undefined" && window.location && window.location.origin
                ? window.location.origin
                : (typeof location !== "undefined" ? location.origin : "");

        const optionsWithOrigin: RequestInit = {
            headers: {
                Origin: origin,
            }
        };

        return this.request<SuggestionResponse>(url, optionsWithOrigin);
    }
}

// Export singleton instance
export const apiClient = new APIClient()

// Export individual functions for convenience
export const getSuggestions = (request: SuggestionRequest) =>
    apiClient.getSuggestions(request)
