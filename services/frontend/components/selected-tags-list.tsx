"use client"

import { translateTag } from "@/lib/translations"
import { X } from "lucide-react"
import { Button } from "@/components/ui/button"

interface SelectedTagsListProps {
    selectedTags: string[]
    onRemoveTag: (index: number) => void
}

export function SelectedTagsList({ selectedTags, onRemoveTag }: SelectedTagsListProps) {
    if (selectedTags.length === 0) {
        return null
    }

    // Group tags by their value and track indices
    const tagGroups: Array<{ tag: string; indices: number[] }> = []
    const tagMap = new Map<string, number[]>()

    selectedTags.forEach((tag, index) => {
        if (!tagMap.has(tag)) {
            tagMap.set(tag, [])
        }
        tagMap.get(tag)!.push(index)
    })

    tagMap.forEach((indices, tag) => {
        tagGroups.push({ tag, indices })
    })

    return (
        <div className="w-full">
            <div className="flex flex-wrap gap-2">
                {tagGroups.map(({ tag, indices }) => {
                    const translatedTag = translateTag(tag)
                    const count = indices.length

                    return (
                        <div
                            key={tag}
                            className="group flex items-center gap-2 bg-white/10 backdrop-blur-md border border-white/20 rounded-full px-3 py-1.5 hover:bg-white/15 hover:border-[#E3B23C]/50 transition-all duration-200 animate-fadeIn"
                        >
                            <span className="text-sm text-[#F8F5F0]">
                                {translatedTag}
                                {count > 1 && (
                                    <span className="ml-1 text-[#E3B23C] font-semibold">
                                        × {count}
                                    </span>
                                )}
                            </span>
                            <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => onRemoveTag(indices[indices.length - 1])}
                                className="h-5 w-5 p-0 rounded-full hover:bg-white/20 text-[#F8F5F0]/60 hover:text-[#F8F5F0] opacity-0 group-hover:opacity-100 transition-opacity duration-200"
                                aria-label={`Удалить ${translatedTag}`}
                            >
                                <X className="h-3 w-3" />
                            </Button>
                        </div>
                    )
                })}
            </div>
        </div>
    )
}

