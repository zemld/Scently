"use client"

import { Button } from "@/components/ui/button"

interface SexSelectorProps {
    value: 'male' | 'unisex' | 'female'
    onChange: (value: 'male' | 'unisex' | 'female') => void
}

export function SexSelector({ value, onChange }: SexSelectorProps) {
    const options = [
        { key: 'male' as const, label: 'Male' },
        { key: 'unisex' as const, label: 'Unisex' },
        { key: 'female' as const, label: 'Female' },
    ]

    return (
        <div className="space-y-2">
            <label className="text-sm font-medium text-[#F8F5F0]/80">Gender</label>
            <div className="flex gap-2">
                {options.map((option) => (
                    <Button
                        key={option.key}
                        variant="ghost"
                        onClick={() => onChange(option.key)}
                        className={`flex-1 py-2 px-4 rounded-xl transition-all duration-200 ${value === option.key
                            ? 'bg-white/10 backdrop-blur-md border border-white/20 text-[#E3B23C] font-semibold'
                            : 'bg-white/5 backdrop-blur-md border border-white/10 text-[#F8F5F0]/60 hover:bg-white/8 hover:text-[#F8F5F0]/80'
                            }`}
                    >
                        {option.label}
                    </Button>
                ))}
            </div>
        </div>
    )
}