<script lang="ts">
    import { fade, scale } from "svelte/transition";
    import { authStore } from "$lib/state/auth.svelte";
    import { toastStore } from "$lib/state/toast.svelte";
    import Button from "./ui/Button.svelte";
    import Input from "./ui/Input.svelte";

    interface Props {
        isOpen: boolean;
        onClose: () => void;
    }

    let { isOpen = $bindable(), onClose }: Props = $props();

    let oldPassword = $state("");
    let newPassword = $state("");
    let confirmPassword = $state("");
    let loading = $state(false);

    async function handleSubmit() {
        if (!oldPassword || !newPassword || !confirmPassword) {
            toastStore.error("Please fill in all fields");
            return;
        }

        if (newPassword !== confirmPassword) {
            toastStore.error("New passwords do not match");
            return;
        }

        if (newPassword.length < 8) {
            toastStore.error("New password must be at least 8 characters long");
            return;
        }

        loading = true;

        try {
            await authStore.changePassword(
                oldPassword,
                newPassword,
                confirmPassword,
            );
            toastStore.success("Password changed successfully!");
            oldPassword = "";
            newPassword = "";
            confirmPassword = "";

            // Close after delay
            setTimeout(() => {
                if (isOpen) handleClose();
            }, 2000);
        } catch (e: any) {
            toastStore.handleApiMessage(e, "Failed to change password");
        } finally {
            loading = false;
        }
    }

    function handleClose() {
        isOpen = false;
        onClose?.();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape" && isOpen) {
            handleClose();
        }
    }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        transition:fade={{ duration: 200 }}
        class="fixed inset-0 z-100 flex items-center justify-center bg-black/80 p-4 backdrop-blur-md"
        onclick={handleClose}
    >
        <div
            transition:scale={{ duration: 300, start: 0.95 }}
            class="relative w-full max-w-lg flex flex-col overflow-hidden rounded-3xl border border-white/10 bg-[#0a0a0a]/90 shadow-2xl"
            onclick={(e) => e.stopPropagation()}
        >
            <!-- Header -->
            <div
                class="flex items-center justify-between border-b border-white/5 p-6 bg-white/5"
            >
                <div>
                    <h2 class="text-xl font-bold text-white tracking-tight">
                        Change Password
                    </h2>
                    <p class="text-xs opacity-80">
                        You'll be logged out from other devices after password
                        change.
                    </p>
                </div>
                <button
                    aria-label="Close"
                    onclick={handleClose}
                    class="rounded-xl p-2 text-white/30 transition-all hover:bg-white/5 hover:text-white active:scale-95"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        ><path d="M18 6 6 18" /><path d="m6 6 12 12" /></svg
                    >
                </button>
            </div>

            <!-- Content -->
            <form
                onsubmit={(e) => {
                    e.preventDefault();
                    handleSubmit();
                }}
                class="p-8 space-y-6"
            >
                <div class="space-y-4">
                    <Input
                        type="password"
                        label="Current Password"
                        placeholder="••••••••"
                        bind:value={oldPassword}
                        disabled={loading}
                        sanitize="password"
                    />

                    <div class="h-px w-full bg-white/5 my-2"></div>

                    <Input
                        type="password"
                        label="New Password"
                        placeholder="••••••••"
                        bind:value={newPassword}
                        disabled={loading}
                        sanitize="password"
                    />

                    <Input
                        type="password"
                        label="Confirm New Password"
                        placeholder="••••••••"
                        bind:value={confirmPassword}
                        disabled={loading}
                        sanitize="password"
                    />
                </div>

                <div class="pt-4 flex flex-col-reverse sm:flex-row gap-3">
                    <Button
                        variant="ghost"
                        class="w-full"
                        onclick={handleClose}
                        disabled={loading}
                    >
                        Cancel
                    </Button>
                    <Button
                        type="submit"
                        variant="primary"
                        class="w-full"
                        {loading}
                    >
                        Update Password
                    </Button>
                </div>
            </form>

            <!-- Footer Decor -->
            <div
                class="h-1 w-full bg-linear-to-r from-brand-primary/50 via-brand-secondary/50 to-brand-primary/50 opacity-30"
            ></div>
        </div>
    </div>
{/if}
