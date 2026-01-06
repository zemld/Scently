"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { usePerfumeAPI } from "@/hooks/use-perfume-api"
import { useSexState } from "@/hooks/use-sex-state"
import { Perfume } from "@/lib/api"
import { PerfumeModal } from "@/components/perfume-modal"
import { SexSelector } from "@/components/sex-selector"
import { RecommendationTypeSwitcher } from "@/components/recommendation-type-switcher"
import { TagCloud } from "@/components/tag-cloud"
import { SelectedTagsList } from "@/components/selected-tags-list"
import { translateSex, translateFamily } from "@/lib/translations"

export default function TagsRecommendationPage() {
    const [selectedTags, setSelectedTags] = useState<string[]>([])
    const [sex, setSex] = useSexState('unisex')
    const [recommendations, setRecommendations] = useState<Perfume[]>([])
    const [selectedPerfume, setSelectedPerfume] = useState<Perfume | null>(null)
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [formMessage, setFormMessage] = useState<string | null>(null)

    const { loading, error, getSuggestionsByTags } = usePerfumeAPI()

    const handleTagClick = (tag: string) => {
        setSelectedTags(prev => [...prev, tag])
        setFormMessage(null)
    }

    const handleRemoveTag = (indexToRemove: number) => {
        setSelectedTags(prev => {
            const newTags = [...prev]
            newTags.splice(indexToRemove, 1)
            return newTags
        })
    }

    const handleClearAll = () => {
        setSelectedTags([])
        setFormMessage(null)
    }

    const handleFindMatches = async () => {
        if (selectedTags.length === 0) {
            setFormMessage("Пожалуйста, выберите хотя бы одно ощущение")
            return
        }

        try {
            setFormMessage(null)
            const results = await getSuggestionsByTags({
                tags: selectedTags,
                sex
            })

            setRecommendations(results)
            setTimeout(() => {
                document.getElementById("recommendations-section")?.scrollIntoView({ behavior: "smooth" })
            }, 100)
        } catch (error) {
            console.error("Failed to get recommendations:", error)
            setFormMessage("Не удалось получить рекомендации. Попробуйте ещё раз.")
        }
    }

    const handlePerfumeClick = (perfume: Perfume) => {
        setSelectedPerfume(perfume)
        setIsModalOpen(true)
    }

    const handleCloseModal = () => {
        setIsModalOpen(false)
        setTimeout(() => setSelectedPerfume(null), 300)
    }

    const defaultPerfumeImage = "/luxury-perfume-bottle-amber-gold.jpg"

    return (
        <div className="min-h-screen bg-gradient-to-br from-[#1C1B1A] via-[#2A1F1E] to-[#5E4B56] text-[#F8F5F0] overflow-x-hidden">
            {/* Hero Section */}
            <section className="min-h-screen flex flex-col items-center justify-center px-4 py-20 relative">
                <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-[#E3B23C] rounded-full blur-[120px] opacity-20 animate-pulse" />
                <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-[#C38E70] rounded-full blur-[120px] opacity-20 animate-pulse delay-1000" />

                <div className="relative z-10 w-full max-w-4xl">
                    <h1 className="text-6xl md:text-7xl font-bold bg-gradient-to-r from-[#C38E70] to-[#E3B23C] bg-clip-text text-transparent text-center mb-4 animate-fadeIn">
                        Рекомендации по ощущениям
                    </h1>
                    <p className="text-xl md:text-2xl text-[#F8F5F0]/80 text-center mb-8 animate-fadeIn delay-200">
                        Выберите ощущения, которые описывают желаемый аромат
                    </p>

                    <RecommendationTypeSwitcher />

                    {/* Input Section */}
                    <div className="w-full bg-white/10 backdrop-blur-xl border border-white/20 rounded-3xl p-8 md:p-12 shadow-[0_8px_32px_rgba(227,178,60,0.2)] animate-expandDown">
                        <h2 className="text-2xl md:text-3xl font-bold text-center mb-8 bg-gradient-to-r from-[#C38E70] to-[#E3B23C] bg-clip-text text-transparent">
                            Выберите ощущения
                        </h2>

                        <div className="space-y-6">
                            <TagCloud selectedTags={selectedTags} onTagClick={handleTagClick} />

                            {selectedTags.length > 0 && (
                                <div className="space-y-3">
                                    <div className="flex items-center justify-between">
                                        <h3 className="text-sm font-medium text-[#F8F5F0]/80">Выбранные ощущения:</h3>
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={handleClearAll}
                                            className="bg-white/5 backdrop-blur-md border-white/20 text-[#F8F5F0]/80 hover:bg-white/10 hover:text-[#F8F5F0] hover:border-white/30 transition-all duration-200 text-xs"
                                        >
                                            Очистить все
                                        </Button>
                                    </div>
                                    <SelectedTagsList selectedTags={selectedTags} onRemoveTag={handleRemoveTag} />
                                </div>
                            )}

                            <SexSelector value={sex} onChange={setSex} />

                            <Button
                                onClick={handleFindMatches}
                                disabled={loading || selectedTags.length === 0}
                                className="w-full py-6 text-lg font-semibold bg-gradient-to-r from-[#C38E70] to-[#E3B23C] hover:from-[#C38E70]/90 hover:to-[#E3B23C]/90 hover:shadow-[0_0_30px_rgba(227,178,60,0.6)] transition-all duration-300 rounded-xl text-white disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {loading ? "Подбор парфюмов..." : "Подобрать парфюмы"}
                            </Button>

                            {formMessage && (
                                <div className="text-[#F8F5F0]/80 text-sm text-center mt-2">
                                    {formMessage}
                                </div>
                            )}

                            {error && (
                                <div className="text-red-400 text-sm text-center mt-2">
                                    {error}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </section>

            {/* Recommendations Section */}
            {recommendations.length > 0 && (
                <section
                    id="recommendations-section"
                    className="w-full px-4 py-20 pb-32 animate-expandDown"
                >
                    <div className="w-full max-w-6xl mx-auto">
                        <h2 className="text-3xl md:text-4xl font-bold text-center mb-12 bg-gradient-to-r from-[#C38E70] to-[#E3B23C] bg-clip-text text-transparent">
                            Ваши идеальные совпадения
                        </h2>

                        <div className="grid md:grid-cols-2 gap-6">
                            {recommendations.map((perfume, index) => (
                                <div
                                    key={perfume.id || index}
                                    onClick={() => handlePerfumeClick(perfume)}
                                    className="group bg-white/10 backdrop-blur-xl border border-white/20 rounded-2xl p-6 hover:bg-white/15 hover:border-[#E3B23C]/50 hover:shadow-[0_8px_32px_rgba(227,178,60,0.3)] hover:scale-105 transition-all duration-300 cursor-pointer animate-fadeIn"
                                    style={{ animationDelay: `${index * 100}ms` }}
                                >
                                    <div className="aspect-square mb-4 rounded-xl overflow-hidden bg-white">
                                        <img
                                            src={perfume.image_url || defaultPerfumeImage}
                                            alt={perfume.name}
                                            className={`w-full h-full group-hover:scale-105 transition-transform duration-300 ${perfume.image_url
                                                ? "object-contain p-4"
                                                : "object-cover"
                                                }`}
                                        />
                                    </div>
                                    <h3 className="text-xl font-bold text-[#F8F5F0] mb-1">{perfume.name}</h3>
                                    <p className="text-sm text-[#C38E70] mb-2">{perfume.brand}</p>
                                    {perfume.properties?.perfume_type && (
                                        <p className="text-sm text-[#F8F5F0]/70">{perfume.properties.perfume_type}</p>
                                    )}
                                    {perfume.similarity_score !== undefined && perfume.similarity_score !== null && perfume.similarity_score > 0 && (
                                        <div className="mt-3 space-y-2">
                                            <div className="flex items-center justify-between text-xs">
                                                <span className="text-[#F8F5F0]/70">Схожесть</span>
                                                <span className="text-[#E3B23C] font-semibold">
                                                    {perfume.similarity_score > 1
                                                        ? `${Math.round(perfume.similarity_score)}%`
                                                        : `${Math.round(perfume.similarity_score * 100)}%`}
                                                </span>
                                            </div>
                                            <div className="relative h-2 w-full overflow-hidden rounded-full bg-white/10">
                                                <div
                                                    className="h-full bg-gradient-to-r from-[#C38E70] to-[#E3B23C] transition-all duration-500 ease-out"
                                                    style={{
                                                        width: `${perfume.similarity_score > 1 ? perfume.similarity_score : perfume.similarity_score * 100}%`
                                                    }}
                                                />
                                            </div>
                                        </div>
                                    )}
                                    {perfume.sex && (
                                        <div className="mt-2">
                                            <span className="text-xs bg-white/10 backdrop-blur-md border border-white/20 text-[#E3B23C] px-2 py-1 rounded-full">
                                                {translateSex(perfume.sex)}
                                            </span>
                                        </div>
                                    )}
                                    {perfume.properties?.family && perfume.properties.family.length > 0 && (
                                        <div className="mt-2 flex flex-wrap gap-1">
                                            {perfume.properties.family.slice(0, 3).map((family: string, familyIndex: number) => (
                                                <span
                                                    key={familyIndex}
                                                    className="text-xs bg-[#C38E70]/20 text-[#C38E70] px-2 py-1 rounded-full"
                                                >
                                                    {translateFamily(family)}
                                                </span>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                </section>
            )}

            {/* Perfume Modal */}
            <PerfumeModal perfume={selectedPerfume} isOpen={isModalOpen} onClose={handleCloseModal} />

            {/* Extra space for smooth animations */}
            <div className="h-32"></div>
        </div>
    )
}

