<script lang="ts">
    interface Props {
        opacity?: number; // 0 to 100 scale
        size?: number;    // pixel cell size
        color?: string;   // stroke line color
        class?: string;   // additional wrapper classes
    }

    let {
        opacity = 5,
        size = 20,
        color = "currentColor",
        class: className = "",
    }: Props = $props();

    // Unique pattern ID to prevent SVG defs collision
    const patternId = `grid-pattern-${Math.random().toString(36).substring(2, 9)}`;
</script>

<div
    class="absolute inset-0 pointer-events-none overflow-hidden select-none -z-10 {className}"
    style:opacity={opacity / 100}
>
    <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
        <defs>
            <pattern
                id={patternId}
                width={size}
                height={size}
                patternUnits="userSpaceOnUse"
            >
                <path
                    d="M {size} 0 L 0 0 0 {size}"
                    fill="none"
                    stroke={color}
                    stroke-width="1"
                />
            </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#{patternId})" />
    </svg>
</div>
