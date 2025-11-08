"use client"

import { Button } from "@/components/ui/button"
import { Perfume } from "@/lib/api"

interface PerfumeModalProps {
    perfume: Perfume | null
    isOpen: boolean
    onClose: () => void
}

export function PerfumeModal({ perfume, isOpen, onClose }: PerfumeModalProps) {
    if (!isOpen || !perfume) return null

    const upperNotes = perfume.properties?.upper_notes || []
    const coreNotes = perfume.properties?.core_notes || []
    const baseNotes = perfume.properties?.base_notes || []

    const hasNotes = upperNotes.length > 0 || coreNotes.length > 0 || baseNotes.length > 0

    // Get unique families (remove duplicates)
    const uniqueFamilies = perfume.properties?.family || []
    const uniqueFamiliesList = uniqueFamilies.length > 0 ? [...new Set(uniqueFamilies)] : []

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 animate-fadeIn" onClick={onClose}>
            {/* Backdrop */}
            <div className="absolute inset-0 bg-black/60 backdrop-blur-md" />

            {/* Modal */}
            <div
                className="relative w-full max-w-4xl max-h-[90vh] overflow-y-auto bg-white/20 backdrop-blur-xl border border-white/30 rounded-3xl shadow-[0_8px_32px_rgba(227,178,60,0.3)] animate-scaleIn"
                onClick={(e) => e.stopPropagation()}
            >

                <div className="p-6 md:p-8">
                    {/* Header Area */}
                    <div className="grid md:grid-cols-2 gap-8 mb-8">
                        {/* Image */}
                        <div className="aspect-square rounded-2xl overflow-hidden bg-white flex items-center justify-center p-4">
                            <img
                                src={perfume.image_url || "/luxury-perfume-bottle-amber-gold.jpg"}
                                alt={perfume.name}
                                className="w-full h-full object-contain"
                            />
                        </div>

                        {/* Info */}
                        <div className="flex flex-col justify-center space-y-4">
                            <p className="text-sm text-[#C38E70] uppercase tracking-wider">{perfume.brand}</p>
                            <h2 className="text-4xl md:text-5xl font-bold text-[#F8F5F0]">{perfume.name}</h2>
                            {perfume.sex && (
                                <div className="inline-flex items-center gap-2 px-4 py-2 bg-white/15 backdrop-blur-md border border-white/30 rounded-full w-fit">
                                    <span className="text-sm text-[#E3B23C]">{perfume.sex}</span>
                                </div>
                            )}
                            {uniqueFamiliesList.length > 0 && (
                                <div className="flex flex-wrap gap-2">
                                    {uniqueFamiliesList.map((family, index) => (
                                        <div key={index} className="inline-flex items-center gap-2 px-4 py-2 bg-white/15 backdrop-blur-md border border-white/30 rounded-full">
                                            <span className="text-sm text-[#E3B23C]">{family}</span>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Notes Section - show if we have any notes */}
                    {hasNotes && (
                        <div className="mb-8">
                            <h3 className="text-2xl font-bold text-[#F8F5F0] mb-6">Ноты аромата</h3>
                            <div className="grid md:grid-cols-3 gap-6">
                                {/* Top Notes */}
                                {upperNotes.length > 0 && (
                                    <div className="space-y-3">
                                        <h4 className="text-sm font-semibold text-[#C38E70] uppercase tracking-wider">Верхние ноты</h4>
                                        <div className="flex flex-wrap gap-2">
                                            {upperNotes.map((note, index) => (
                                                <span
                                                    key={index}
                                                    className="px-3 py-1.5 text-sm bg-white/15 backdrop-blur-md border border-white/30 rounded-full text-[#F8F5F0] hover:bg-white/20 hover:border-[#E3B23C]/50 transition-all duration-300"
                                                >
                                                    {note}
                                                </span>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {/* Core Notes */}
                                {coreNotes.length > 0 && (
                                    <div className="space-y-3">
                                        <h4 className="text-sm font-semibold text-[#C38E70] uppercase tracking-wider">Ноты сердца</h4>
                                        <div className="flex flex-wrap gap-2">
                                            {coreNotes.map((note, index) => (
                                                <span
                                                    key={index}
                                                    className="px-3 py-1.5 text-sm bg-white/15 backdrop-blur-md border border-white/30 rounded-full text-[#F8F5F0] hover:bg-white/20 hover:border-[#E3B23C]/50 transition-all duration-300"
                                                >
                                                    {note}
                                                </span>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {/* Base Notes */}
                                {baseNotes.length > 0 && (
                                    <div className="space-y-3">
                                        <h4 className="text-sm font-semibold text-[#C38E70] uppercase tracking-wider">Базовые ноты</h4>
                                        <div className="flex flex-wrap gap-2">
                                            {baseNotes.map((note, index) => (
                                                <span
                                                    key={index}
                                                    className="px-3 py-1.5 text-sm bg-white/15 backdrop-blur-md border border-white/30 rounded-full text-[#F8F5F0] hover:bg-white/20 hover:border-[#E3B23C]/50 transition-all duration-300"
                                                >
                                                    {note}
                                                </span>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    {/* Available Shops and Variants - show if we have shops data */}
                    {perfume.shops && perfume.shops.length > 0 && (
                        <div className="mb-6">
                            <h3 className="text-2xl font-bold text-[#F8F5F0] mb-6">Где купить</h3>
                            <div className="space-y-6">
                                {perfume.shops.map((shop, shopIndex) => (
                                    <div key={shopIndex} className="bg-white/5 backdrop-blur-md border border-white/20 rounded-xl p-4">
                                        <div className="flex items-center justify-between mb-3">
                                            <h4 className="text-lg font-semibold text-[#F8F5F0]">{shop.shop_name}</h4>
                                            {shop.domain && (
                                                <a
                                                    href={`https://${shop.domain}`}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="text-sm text-[#C38E70] hover:text-[#E3B23C] transition-colors"
                                                >
                                                    {shop.domain}
                                                </a>
                                            )}
                                        </div>
                                        {shop.variants && shop.variants.length > 0 && (
                                            <div className="flex flex-wrap gap-3">
                                                {shop.variants.map((variant, variantIndex) => (
                                                    <a
                                                        key={variantIndex}
                                                        href={variant.link}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="px-4 py-2 bg-white/15 backdrop-blur-md border border-white/30 rounded-lg text-[#F8F5F0] hover:bg-white/20 hover:border-[#E3B23C]/50 hover:shadow-[0_0_20px_rgba(227,178,60,0.5)] hover:scale-105 transition-all duration-300"
                                                    >
                                                        <div className="flex flex-col">
                                                            <span className="font-semibold">{variant.volume} мл</span>
                                                            {variant.price > 0 && (
                                                                <span className="text-xs text-[#C38E70]">{variant.price} ₽</span>
                                                            )}
                                                        </div>
                                                    </a>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}

                    {/* Footer */}
                    <div className="flex justify-center pt-6 border-t border-white/20">
                        <Button
                            onClick={onClose}
                            className="px-8 py-3 bg-gradient-to-r from-[#C38E70] to-[#E3B23C] hover:from-[#C38E70]/90 hover:to-[#E3B23C]/90 hover:shadow-[0_0_30px_rgba(227,178,60,0.6)] transition-all duration-300 rounded-xl text-white font-semibold"
                        >
                            Закрыть
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    )
}
