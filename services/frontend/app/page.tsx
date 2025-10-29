"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Switch } from "@/components/ui/switch"
import { usePerfumeAPI } from "@/hooks/use-perfume-api"
import { Perfume } from "@/lib/api"
import { PerfumeModal } from "@/components/perfume-modal"
import { SexSelector } from "@/components/sex-selector"

export default function ScentlyLanding() {
  const [showInputSection, setShowInputSection] = useState(false)
  const [showRecommendations, setShowRecommendations] = useState(false)
  const [brand, setBrand] = useState("")
  const [name, setName] = useState("")
  const [sex, setSex] = useState<'male' | 'unisex' | 'female'>('unisex')
  const [aiMode, setAiMode] = useState(false)
  const [recommendations, setRecommendations] = useState<Perfume[]>([])
  const [selectedPerfume, setSelectedPerfume] = useState<Perfume | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [formMessage, setFormMessage] = useState<string | null>(null)

  const { loading, error, getSuggestions } = usePerfumeAPI()

  const handleTryNow = () => {
    setShowInputSection(true)
    setTimeout(() => {
      document.getElementById("input-section")?.scrollIntoView({ behavior: "smooth" })
    }, 100)
  }

  const handleFindMatches = async () => {
    if (!brand.trim() || !name.trim()) {
      setFormMessage("Пожалуйста, укажите бренд и название")
      return
    }

    try {
      setFormMessage(null)
      const results = await getSuggestions({
        brand,
        name,
        use_ai: aiMode,
        sex
      })

      setRecommendations(results)
      setShowRecommendations(true)
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

  // Fallback images for recommendations
  const fallbackImages = [
    "/luxury-perfume-bottle-dark-wood.jpg",
    "/luxury-perfume-bottle-silver-elegant.jpg",
    "/luxury-perfume-bottle-amber-gold.jpg",
    "/luxury-perfume-bottle-red-crystal.jpg"
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#1C1B1A] via-[#2A1F1E] to-[#5E4B56] text-[#F8F5F0] overflow-x-hidden">
      {/* Hero Section */}
      <section className="min-h-screen flex flex-col items-center justify-center px-4 relative">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-[#E3B23C] rounded-full blur-[120px] opacity-20 animate-pulse" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-[#C38E70] rounded-full blur-[120px] opacity-20 animate-pulse delay-1000" />

        <div className="relative z-10 text-center space-y-8">
          <h1 className="text-7xl md:text-8xl font-bold bg-gradient-to-r from-[#C38E70] to-[#E3B23C] bg-clip-text text-transparent animate-fadeIn">
            Scently
          </h1>
          <p className="text-2xl md:text-3xl text-[#F8F5F0]/80 animate-fadeIn delay-200">
            Найдите аромат, созданный именно для вас.
          </p>
          <Button
            onClick={handleTryNow}
            className="mt-8 px-12 py-6 text-lg font-semibold bg-white/10 backdrop-blur-xl border border-white/20 hover:bg-white/20 hover:border-[#E3B23C]/50 hover:shadow-[0_0_30px_rgba(227,178,60,0.5)] transition-all duration-300 rounded-2xl text-[#F8F5F0] animate-fadeIn delay-400"
          >
            Попробовать
          </Button>
        </div>
      </section>

      {/* Interactive Input Section */}
      {showInputSection && (
        <section
          id="input-section"
          className="min-h-screen flex items-center justify-center px-4 py-20 animate-expandDown"
        >
          <div className="w-full max-w-2xl bg-white/10 backdrop-blur-xl border border-white/20 rounded-3xl p-8 md:p-12 shadow-[0_8px_32px_rgba(227,178,60,0.2)]">
            <h2 className="text-3xl md:text-4xl font-bold text-center mb-8 bg-gradient-to-r from-[#C38E70] to-[#E3B23C] bg-clip-text text-transparent">
              Расскажите о ваших предпочтениях
            </h2>

            <div className="space-y-6">
              <div className="grid md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-[#F8F5F0]/80">Бренд</label>
                  <Input
                    placeholder="например, Dior"
                    value={brand}
                    onChange={(e) => setBrand(e.target.value)}
                    className="bg-white/5 backdrop-blur-md border-white/20 focus:border-[#C38E70] focus:ring-[#C38E70]/50 rounded-xl text-[#F8F5F0] placeholder:text-[#F8F5F0]/40"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-[#F8F5F0]/80">Название</label>
                  <Input
                    placeholder="например, Sauvage"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="bg-white/5 backdrop-blur-md border-white/20 focus:border-[#C38E70] focus:ring-[#C38E70]/50 rounded-xl text-[#F8F5F0] placeholder:text-[#F8F5F0]/40"
                  />
                </div>
              </div>

              <SexSelector value={sex} onChange={setSex} />

              <div className="flex items-center justify-between bg-white/5 backdrop-blur-md border border-white/20 rounded-xl p-4">
                <label className="text-sm font-medium text-[#F8F5F0]">Использовать ИИ</label>
                <Switch checked={aiMode} onCheckedChange={setAiMode} className="data-[state=checked]:bg-[#C38E70]" />
              </div>

              <Button
                onClick={handleFindMatches}
                disabled={loading}
                className="w-full py-6 text-lg font-semibold bg-gradient-to-r from-[#C38E70] to-[#E3B23C] hover:from-[#C38E70]/90 hover:to-[#E3B23C]/90 hover:shadow-[0_0_30px_rgba(227,178,60,0.6)] transition-all duration-300 rounded-xl text-white disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? "Поиск совпадений..." : "Найти похожие"}
              </Button>

              {formMessage && (
                <div className="text-[#F8F5F0]/80 text-sm text-center mt-2">
                  {formMessage}
                </div>
              )}
            </div>
          </div>
        </section>
      )}

      {/* Recommendations Section */}
      {showRecommendations && (
        <section
          id="recommendations-section"
          className="w-full px-4 py-20 pb-32 animate-expandDown"
        >
          <div className="w-full max-w-6xl mx-auto">
            <h2 className="text-3xl md:text-4xl font-bold text-center mb-12 bg-gradient-to-r from-[#C38E70] to-[#E3B23C] bg-clip-text text-transparent">
              Ваши идеальные совпадения
            </h2>

            {recommendations.length > 0 ? (
              <div className="grid md:grid-cols-2 gap-6">
                {recommendations.map((perfume, index) => (
                  <div
                    key={perfume.id || index}
                    onClick={() => handlePerfumeClick(perfume)}
                    className="group bg-white/10 backdrop-blur-xl border border-white/20 rounded-2xl p-6 hover:bg-white/15 hover:border-[#E3B23C]/50 hover:shadow-[0_8px_32px_rgba(227,178,60,0.3)] hover:scale-105 transition-all duration-300 cursor-pointer animate-fadeIn"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <div className="aspect-square mb-4 rounded-xl overflow-hidden bg-white/5">
                      <img
                        src={perfume.image || fallbackImages[index % fallbackImages.length]}
                        alt={perfume.name}
                        className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                      />
                    </div>
                    <h3 className="text-xl font-bold text-[#F8F5F0] mb-1">{perfume.name}</h3>
                    <p className="text-sm text-[#C38E70] mb-2">{perfume.brand}</p>
                    <p className="text-sm text-[#F8F5F0]/70">{perfume.description}</p>
                    {perfume.sex && (
                      <div className="mt-2">
                        <span className="text-xs bg-white/10 backdrop-blur-md border border-white/20 text-[#E3B23C] px-2 py-1 rounded-full">
                          {perfume.sex}
                        </span>
                      </div>
                    )}
                    {perfume.family && (
                      <div className="mt-2 flex flex-wrap gap-1">
                        {perfume.family.slice(0, 3).map((family, familyIndex) => (
                          <span
                            key={familyIndex}
                            className="text-xs bg-[#C38E70]/20 text-[#C38E70] px-2 py-1 rounded-full"
                          >
                            {family}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <p className="text-[#F8F5F0]/60 text-lg">Рекомендации не найдены. Попробуйте изменить предпочтения.</p>
              </div>
            )}
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
