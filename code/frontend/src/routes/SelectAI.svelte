<script lang="ts">
    import { AIs } from "$lib/components/AIs";
    import { onMount, onDestroy, createEventDispatcher } from "svelte";

    interface Props {
        title: string;
        onSelectAI: (ai: string) => void;
        onClose?: () => void;
    }

    let { title, onSelectAI, onClose }: Props = $props();

    const dispatch = createEventDispatcher();

    // lazy loaders returning URLs for matched files
    const modules = import.meta.glob("/src/lib/assets/ai/*.{png,jpg}", {
        as: "url",
    }) as Record<string, () => Promise<string>>;

    function findKey(name: string) {
        return Object.keys(modules).find(
            (p) => p.endsWith(`/${name}.png`) || p.endsWith(`/${name}.jpg`)
        );
    }

    async function loadUrl(name: string) {
        const key = findKey(name);
        if (!key) {
            // return a fallback path (will likely 404 and show browser broken-image icon)
            return `/lib/assets/ai/${name}.png`;
        }
        return await modules[key]();
    }

    // cache resolved URLs to avoid repeated loads
    const cache: Record<string, string> = {};
    async function getUrl(ai: string) {
        if (cache[ai]) return cache[ai];
        const url = await loadUrl(ai);
        cache[ai] = url;
        return url;
    }

    // accessibility & focus management
    let previousActive: Element | null = null;
    let dialogId = `selectai-${Math.random().toString(36).slice(2, 10)}`;

    function close() {
        if (onClose) onClose();
        dispatch("close");
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape" || e.key === "Esc") {
            e.preventDefault();
            close();
        }
    }

    onMount(() => {
        previousActive = document.activeElement;
        // prevent background scroll while modal is open
        const prevOverflow = document.body.style.overflow;
        document.body.style.overflow = "hidden";
        window.addEventListener("keydown", handleKeydown);
        // focus the dialog for accessibility
        const el = document.getElementById(dialogId);
        if (el) (el as HTMLElement).focus();

        return () => {
            document.body.style.overflow = prevOverflow;
            window.removeEventListener("keydown", handleKeydown);
            if (previousActive && (previousActive as HTMLElement).focus) {
                try {
                    (previousActive as HTMLElement).focus();
                } catch (e) {
                    /* ignore */
                }
            }
        };
    });

    onDestroy(() => {
        // onMount's returned cleanup runs automatically, keep guard
    });

    function selectAI(id: string) {
        onSelectAI(id);
        close();
    }
</script>

<div class="selectAIbackdrop" role="presentation">
    <div id={dialogId} class="selectAI" role="dialog">
        <button class="close" aria-label="Close dialog" onclick={close}>
            ✕
        </button>

        <h2 id={`title-${dialogId}`}>{title}</h2>

        <div class="grid">
            {#each AIs as ai}
                <button
                    class="card"
                    type="button"
                    onclick={() => selectAI(ai.id)}
                >
                    <div class="card-inner">
                        {#await getUrl(ai.id) then url}
                            <img
                                src={url}
                                alt={`AI: ${ai.name}`}
                                loading="lazy"
                            />
                        {:catch}
                            <div class="img-fallback" aria-hidden="true">
                                🤖
                            </div>
                        {/await}

                        <div class="info">
                            <h3>{ai.name}</h3>
                            <p>{ai.description}</p>
                        </div>
                    </div>
                </button>
            {/each}
        </div>
    </div>
</div>

<style>
    .selectAIbackdrop {
        position: fixed;
        inset: 0;
        background-color: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;

        .selectAI {
            background: var(--bg);
            color: var(--text);
            border-radius: 10px;
            padding: 20px 24px;
            width: min(1000px, 92%);
            max-height: 90vh;
            overflow: auto;
            box-shadow: 0 12px 30px rgba(0, 0, 0, 0.35);
            position: relative;
            outline: none;

            h2 {
                margin: 0 0 12px 0;
                font-size: 1.5rem;
                text-align: center;
            }

            .close {
                position: absolute;
                right: 12px;
                top: 12px;
                background: transparent;
                border: none;
                color: var(--text);
                font-size: 1.25rem;
                cursor: pointer;
                padding: 6px;
            }

            .grid {
                display: grid;
                grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
                gap: 12px;

                .card {
                    background: var(--bg-accent);
                    color: var(--text);
                    border: none;
                    padding: 0;
                    text-align: left;
                    border-radius: 8px;
                    cursor: pointer;
                    transition:
                        transform 0.12s ease,
                        box-shadow 0.12s ease;
                    display: block;
                    width: 100%;

                    &:focus,
                    &:hover {
                        transform: translateY(-4px);
                        box-shadow: 0 8px 22px rgba(0, 0, 0, 0.18);
                        outline: none;
                    }

                    .card-inner {
                        display: flex;
                        align-items: center;
                        gap: 12px;
                        padding: 12px;

                        img {
                            width: 128px;
                            height: 128px;
                            object-fit: cover;
                            border-radius: 6px;
                            flex-shrink: 0;
                        }

                        .img-fallback {
                            width: 128px;
                            height: 128px;
                            display: inline-grid;
                            place-items: center;
                            border-radius: 6px;

                            color: #666;
                            font-size: 1.2rem;
                        }

                        .info {
                            h3 {
                                margin: 0 0 4px 0;
                                font-size: 1.05rem;
                            }
                            p {
                                margin: 0;
                                color: var(--muted, #9aa);
                                font-size: 0.92rem;
                            }
                        }
                    }
                }
            }
        }
    }

    @media (max-width: 520px) {
        .selectAI {
            padding: 16px;
        }
        img,
        .img-fallback {
            width: 52px;
            height: 52px;
        }
    }
</style>
