"use client"

import { tagsTranslations, translateTag } from "@/lib/translations"
import { cn } from "@/lib/utils"
import { useState } from "react"

interface TagCloudProps {
    selectedTags: string[]
    onTagClick: (tag: string) => void
}

export function TagCloud({ selectedTags, onTagClick }: TagCloudProps) {
    // Get all available tags (English keys)
    const allTags = Object.keys(tagsTranslations)

    // Create a Set for quick lookup of selected tags
    const selectedTagsSet = new Set(selectedTags)
    const [clickedTag, setClickedTag] = useState<string | null>(null)

    const handleTagClick = (tag: string) => {
        setClickedTag(tag)
        onTagClick(tag)
        // Reset animation after it completes
        setTimeout(() => setClickedTag(null), 600)
    }

    return (
        <div className="w-full">
            <div className="flex flex-wrap gap-3 justify-center p-4">
                {allTags.map((tag) => {
                    const isSelected = selectedTagsSet.has(tag)
                    const isClicked = clickedTag === tag
                    const translatedTag = translateTag(tag)

                    return (
                        <button
                            key={tag}
                            onClick={() => handleTagClick(tag)}
                            className={cn(
                                "px-4 py-2 rounded-full text-sm font-medium transition-all duration-300",
                                "backdrop-blur-md border",
                                isClicked && "scale-110 shadow-[0_0_30px_rgba(227,178,60,0.8)]",
                                isSelected
                                    ? "bg-gradient-to-r from-[#C38E70]/30 to-[#E3B23C]/30 border-[#E3B23C]/50 text-[#E3B23C] shadow-[0_0_20px_rgba(227,178,60,0.4)] hover:shadow-[0_0_25px_rgba(227,178,60,0.6)]"
                                    : "bg-white/5 border-white/20 text-[#F8F5F0]/80 hover:bg-white/10 hover:border-white/30 hover:text-[#F8F5F0] hover:shadow-[0_0_15px_rgba(255,255,255,0.2)] hover:scale-105 active:scale-95"
                            )}
                        >
                            {translatedTag}
                        </button>
                    )
                })}
            </div>
        </div>
    )
}

