"use client"

import { useRef, useEffect, useState } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"

export function RecommendationTypeSwitcher() {
    const pathname = usePathname()
    const isTagsPage = pathname === "/tags"

    const linkRefs = useRef<(HTMLAnchorElement | null)[]>([])
    const containerRef = useRef<HTMLDivElement>(null)
    const [sliderStyle, setSliderStyle] = useState({ left: 0, width: 0 })

    const activeIndex = isTagsPage ? 1 : 0

    useEffect(() => {
        const updateSlider = () => {
            const activeLink = linkRefs.current[activeIndex]
            const container = containerRef.current

            if (activeLink && container) {
                const containerRect = container.getBoundingClientRect()
                const linkRect = activeLink.getBoundingClientRect()

                setSliderStyle({
                    left: linkRect.left - containerRect.left,
                    width: linkRect.width
                })
            }
        }

        // Small delay to ensure links are rendered
        const timeoutId = setTimeout(updateSlider, 0)
        window.addEventListener('resize', updateSlider)
        return () => {
            clearTimeout(timeoutId)
            window.removeEventListener('resize', updateSlider)
        }
    }, [activeIndex])

    return (
        <div className="flex justify-center mb-8">
            <div
                ref={containerRef}
                className="relative bg-white/5 backdrop-blur-xl border border-white/20 rounded-2xl p-2 flex gap-4 shadow-[0_8px_32px_rgba(227,178,60,0.1)]"
            >
                {/* Liquid glass background for active item */}
                <div
                    className="absolute top-2 bottom-2 rounded-xl bg-white/10 backdrop-blur-md border border-white/20 transition-all duration-300 ease-out"
                    style={{
                        left: `${sliderStyle.left}px`,
                        width: `${sliderStyle.width}px`,
                        boxShadow: isTagsPage
                            ? "0_0_30px_rgba(227,178,60,0.5), inset 0_1px_0_rgba(255,255,255,0.2)"
                            : "0_0_30px_rgba(195,142,112,0.5), inset 0_1px_0_rgba(255,255,255,0.2)"
                    }}
                />

                <Link
                    ref={(el) => { linkRefs.current[0] = el }}
                    href="/similar"
                    className={cn(
                        "relative z-10 px-6 py-3 rounded-xl font-semibold transition-all duration-300 flex-1 flex items-center justify-center whitespace-nowrap",
                        !isTagsPage
                            ? "text-[#E3B23C]"
                            : "text-[#F8F5F0]/60 hover:text-[#F8F5F0]/80"
                    )}
                >
                    По парфюму
                </Link>

                <Link
                    ref={(el) => { linkRefs.current[1] = el }}
                    href="/tags"
                    className={cn(
                        "relative z-10 px-6 py-3 rounded-xl font-semibold transition-all duration-300 flex-1 flex items-center justify-center whitespace-nowrap",
                        isTagsPage
                            ? "text-[#E3B23C]"
                            : "text-[#F8F5F0]/60 hover:text-[#F8F5F0]/80"
                    )}
                >
                    По ощущениям
                </Link>
            </div>
        </div>
    )
}
