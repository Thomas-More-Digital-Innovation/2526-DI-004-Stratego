<script lang="ts">
    import changelogRaw from "virtual:changelog";

    function parseMarkdown(md: string) {
        if (!md) return "";

        return md
            .trim()
            .replace(
                /^# (.*$)/gim,
                '<h1 class="text-4xl font-extrabold mb-8 bg-gradient-to-r from-white to-white/40 bg-clip-text text-transparent">$1</h1>',
            )
            .replace(
                /^## \[(.*)\] - (.*$)/gim,
                `
                <div class="mt-12 mb-6">
                    <div class="flex items-center gap-4">
                        <h2 class="text-2xl font-bold text-white">$1</h2>
                        <span class="text-sm font-medium text-white/30">$2</span>
                    </div>
                    <div class="mt-2 h-px w-full bg-gradient-to-r from-white/10 to-transparent"></div>
                </div>
            `,
            )
            .replace(
                /^### (.*$)/gim,
                '<h3 class="text-lg font-semibold mt-8 mb-3 text-brand-primary/90 flex items-center gap-2"><span class="h-1.5 w-1.5 rounded-full bg-brand-primary"></span>$1</h3>',
            )
            .replace(
                /^\* \*\*(.*)\*\*(.*)$/gim,
                '<li class="ml-2 mb-2 flex gap-3 text-white/80"><span class="text-white/20">•</span><span><strong class="text-white">$1</strong>$2</span></li>',
            )
            .replace(
                /^\- (.*$)/gim,
                '<li class="ml-2 mb-2 flex gap-3 text-white/80"><span class="text-white/20">•</span><span>$1</span></li>',
            )
            .replace(
                /\*\*(.*)\*\*/gim,
                '<strong class="text-white">$1</strong>',
            )
            .split("\n\n")
            .map((p) =>
                p.trim().startsWith("<")
                    ? p
                    : `<p class="mb-4 text-white/60 leading-relaxed">${p}</p>`,
            )
            .join("\n");
    }

    const html = parseMarkdown(changelogRaw);
</script>

<div class="changelog-content">
    {@html html}
</div>

<style>
    :global(.changelog-content h1:first-child) {
        display: none; /* Hide the main "Changelog" title as it's in the modal header */
    }
</style>
