<script lang="ts">
	import type { PageData } from './$types';
	import { invalidateAll } from '$app/navigation';

	interface BoardSetup {
		id: number;
		user_id: number;
		name: string;
		description: string;
		setup_data: string;
		is_default: boolean;
		created_at: string;
		updated_at: string;
	}

	let { data }: { data: PageData } = $props();

	let setups = $state<BoardSetup[]>(data.setups || []);
	let error = $state(data.error || '');
	let showCreateModal = $state(false);
	let editingSetup = $state<BoardSetup | null>(null);

	// Form state
	let setupName = $state('');
	let setupDescription = $state('');
	let setupData = $state('');
	let isDefault = $state(false);

	function openCreateModal() {
		setupName = '';
		setupDescription = '';
		setupData = '{}';
		isDefault = false;
		editingSetup = null;
		showCreateModal = true;
	}

	function openEditModal(setup: BoardSetup) {
		setupName = setup.name;
		setupDescription = setup.description;
		setupData = setup.setup_data;
		isDefault = setup.is_default;
		editingSetup = setup;
		showCreateModal = true;
	}

	async function saveSetup() {
		if (!setupName.trim()) {
			alert('Setup name is required');
			return;
		}

		try {
			const url = editingSetup
				? `http://localhost:8080/api/board-setups?id=${editingSetup.id}`
				: `http://localhost:8080/api/board-setups`;

			const method = editingSetup ? 'PUT' : 'POST';

			const response = await fetch(url, {
				method,
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({
					name: setupName,
					description: setupDescription,
					setup_data: setupData,
					is_default: isDefault
				})
			});

			if (!response.ok) throw new Error('Failed to save setup');

			showCreateModal = false;
			await invalidateAll(); // Reload server data
		} catch (e) {
			alert('Failed to save: ' + e);
		}
	}

	async function deleteSetup(setupId: number) {
		if (!confirm('Delete this setup?')) return;

		try {
			const response = await fetch(`http://localhost:8080/api/board-setups?id=${setupId}`, {
				method: 'DELETE',
				credentials: 'include'
			});

			if (!response.ok) throw new Error('Failed to delete');

			await invalidateAll(); // Reload server data
		} catch (e) {
			alert('Failed to delete: ' + e);
		}
	}
</script>

<div class="container">
	<div class="header">
		<h1>Board Setups</h1>
		<button class="btn-primary" onclick={openCreateModal}> + Create New Setup </button>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{:else if setups.length === 0}
		<div class="empty">
			<p>No board setups yet. Create your first setup!</p>
		</div>
	{:else}
		<div class="grid">
			{#each setups as setup}
				<div class="card">
					<div class="card-header">
						<h3>{setup.name}</h3>
						{#if setup.is_default}
							<span class="badge">Default</span>
						{/if}
					</div>
					{#if setup.description}
						<p class="description">{setup.description}</p>
					{/if}
					<div class="card-actions">
						<button class="btn-edit" onclick={() => openEditModal(setup)}> Edit </button>
						<button class="btn-delete" onclick={() => deleteSetup(setup.id)}>
							Delete
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

{#if showCreateModal}
	<div class="modal-backdrop" onclick={() => (showCreateModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>{editingSetup ? 'Edit' : 'Create'} Board Setup</h2>

			<div class="field">
				<label for="name">Name</label>
				<input id="name" type="text" bind:value={setupName} placeholder="My Setup" />
			</div>

			<div class="field">
				<label for="description">Description</label>
				<textarea
					id="description"
					bind:value={setupDescription}
					placeholder="Optional description"
					rows="3"
				></textarea>
			</div>

			<div class="field">
				<label for="data">Setup Data (JSON)</label>
				<textarea
					id="data"
					bind:value={setupData}
					placeholder='Enter JSON data'
					rows="6"
				></textarea>
			</div>

			<div class="field-checkbox">
				<input id="default" type="checkbox" bind:checked={isDefault} />
				<label for="default">Set as default</label>
			</div>

			<div class="modal-actions">
				<button class="btn-secondary" onclick={() => (showCreateModal = false)}>
					Cancel
				</button>
				<button class="btn-primary" onclick={saveSetup}>Save</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.container {
		padding: 20px;
		max-width: 1200px;
		margin: 0 auto;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
	}

	h1 {
		margin: 0;
		color: var(--text);
	}

	.error {
		background: #f44;
		color: white;
		padding: 16px;
		border-radius: 8px;
		margin-bottom: 20px;
	}

	.empty {
		text-align: center;
		padding: 60px 20px;
		color: var(--muted, #999);
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 20px;
	}

	.card {
		background: var(--bg-accent);
		padding: 20px;
		border-radius: 10px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 12px;
	}

	.card h3 {
		margin: 0;
		color: var(--text);
	}

	.badge {
		background: #66f;
		color: white;
		padding: 4px 10px;
		border-radius: 12px;
		font-size: 0.85rem;
		font-weight: 600;
	}

	.description {
		color: var(--muted, #999);
		margin: 8px 0;
		font-size: 0.95rem;
	}

	.card-actions {
		display: flex;
		gap: 10px;
		margin-top: 16px;
	}

	.btn-primary {
		background: #66f;
		color: white;
		border: none;
		padding: 10px 20px;
		border-radius: 6px;
		cursor: pointer;
		font-weight: 600;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #55d;
	}

	.btn-secondary {
		background: transparent;
		color: var(--text);
		border: 1px solid #444;
		padding: 10px 20px;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-secondary:hover {
		border-color: #66f;
		color: #66f;
	}

	.btn-edit {
		flex: 1;
		background: #4a4;
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: 6px;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-edit:hover {
		background: #393;
	}

	.btn-delete {
		flex: 1;
		background: #d44;
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: 6px;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-delete:hover {
		background: #c33;
	}

	.modal-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.6);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: var(--bg);
		padding: 30px;
		border-radius: 12px;
		max-width: 500px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal h2 {
		margin: 0 0 20px 0;
		color: var(--text);
	}

	.field {
		margin-bottom: 16px;
	}

	.field label {
		display: block;
		margin-bottom: 6px;
		color: var(--text);
		font-weight: 500;
	}

	.field input,
	.field textarea {
		width: 100%;
		padding: 10px;
		border: 1px solid #444;
		border-radius: 6px;
		background: var(--bg-accent);
		color: var(--text);
		font-size: 1rem;
		font-family: inherit;
	}

	.field textarea {
		resize: vertical;
	}

	.field-checkbox {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-bottom: 20px;
	}

	.field-checkbox input {
		width: auto;
	}

	.modal-actions {
		display: flex;
		gap: 12px;
		justify-content: flex-end;
	}
</style>
