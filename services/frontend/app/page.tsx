"use client"

import Link from "next/link"
import { Button } from "@/components/ui/button"

export default function ScentlyLanding() {
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
          <Link href="/similar">
            <Button
              className="mt-8 px-12 py-6 text-lg font-semibold bg-white/10 backdrop-blur-xl border border-white/20 hover:bg-white/20 hover:border-[#E3B23C]/50 hover:shadow-[0_0_30px_rgba(227,178,60,0.5)] transition-all duration-300 rounded-2xl text-[#F8F5F0] animate-fadeIn delay-400"
            >
              Попробовать
            </Button>
          </Link>
        </div>
      </section>

      {/* Extra space for smooth animations */}
      <div className="h-32"></div>
    </div>
  )
}
