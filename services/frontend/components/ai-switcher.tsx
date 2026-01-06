"use client"

import { useRef, useEffect, useState } from "react"
import { cn } from "@/lib/utils"

interface AISwitcherProps {
    value: boolean
    onChange: (value: boolean) => void
}

export function AISwitcher({ value, onChange }: AISwitcherProps) {
    const buttonRefs = useRef<(HTMLButtonElement | null)[]>([])
    const containerRef = useRef<HTMLDivElement>(null)
    const [sliderStyle, setSliderStyle] = useState({ left: 0, width: 0 })

    const activeIndex = value ? 1 : 0

    useEffect(() => {
        const updateSlider = () => {
            const activeButton = buttonRefs.current[activeIndex]
            const container = containerRef.current

            if (activeButton && container) {
                const containerRect = container.getBoundingClientRect()
                const buttonRect = activeButton.getBoundingClientRect()

                setSliderStyle({
                    left: buttonRect.left - containerRect.left,
                    width: buttonRect.width
                })
            }
        }

        const timeoutId = setTimeout(updateSlider, 0)
        window.addEventListener('resize', updateSlider)
        return () => {
            clearTimeout(timeoutId)
            window.removeEventListener('resize', updateSlider)
        }
    }, [activeIndex])

    return (
        <div className="space-y-2">
            <label className="text-sm font-medium text-[#F8F5F0]/80">Использовать ИИ</label>
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
                        boxShadow: value
                            ? "0_0_30px_rgba(227,178,60,0.5), inset 0_1px_0_rgba(255,255,255,0.2)"
                            : "0_0_30px_rgba(195,142,112,0.5), inset 0_1px_0_rgba(255,255,255,0.2)"
                    }}
                />

                <button
                    ref={(el) => { buttonRefs.current[0] = el }}
                    onClick={() => onChange(false)}
                    className={cn(
                        "relative z-10 px-6 py-3 rounded-xl font-semibold transition-all duration-300 flex-1 flex items-center justify-center whitespace-nowrap",
                        !value
                            ? "text-[#E3B23C]"
                            : "text-[#F8F5F0]/60 hover:text-[#F8F5F0]/80"
                    )}
                >
                    Выкл
                </button>

                <button
                    ref={(el) => { buttonRefs.current[1] = el }}
                    onClick={() => onChange(true)}
                    className={cn(
                        "relative z-10 px-6 py-3 rounded-xl font-semibold transition-all duration-300 flex-1 flex items-center justify-center whitespace-nowrap",
                        value
                            ? "text-[#E3B23C]"
                            : "text-[#F8F5F0]/60 hover:text-[#F8F5F0]/80"
                    )}
                >
                    Вкл
                </button>
            </div>
        </div>
    )
}
