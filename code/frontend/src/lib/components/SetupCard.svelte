<script lang="ts">
	import { PIECE_INVENTORY, type BoardSetup } from '$lib/types/boardSetup';

	interface Props {
		setup: BoardSetup;
		onEdit: (setup: BoardSetup) => void;
		onDelete: (id: number) => void;
	}

	let { setup, onEdit, onDelete }: Props = $props();

	// Preview first few pieces
	function getPreviewPieces(): string[] {
		const data = setup.setup_data;
		if (data.length >= 10) {
			return data.slice(0, 10).split('');
		}
		return Array(10).fill('.');
	}
</script>

<div class="setup-card">
	<div class="card-header">
		<h3>{setup.name}</h3>
		{#if setup.is_default}
			<span class="badge">Default</span>
		{/if}
	</div>

	{#if setup.description}
		<p class="description">{setup.description}</p>
	{/if}

	<div class="preview">
		{#each getPreviewPieces() as piece}
			<div class="preview-cell">
				{#if piece !== '.' && PIECE_INVENTORY[piece]}
					<div class="mini-piece">
						<span class="icon">{PIECE_INVENTORY[piece].icon}</span>
						<span class="rank">{PIECE_INVENTORY[piece].rank}</span>
					</div>
				{/if}
			</div>
		{/each}
	</div>

	<div class="card-actions">
		<button class="btn-edit" onclick={() => onEdit(setup)}>Edit</button>
		<button class="btn-delete" onclick={() => onDelete(setup.id)}>Delete</button>
	</div>
</div>

<style>
	.setup-card {
		background: var(--bg-accent);
		padding: 20px;
		border-radius: 10px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		transition: transform 0.2s;

		&:hover {
			transform: translateY(-2px);
		}
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 12px;

		h3 {
			margin: 0;
			color: var(--text);
			font-size: 1.1rem;
		}

		.badge {
			background: #66f;
			color: white;
			padding: 4px 10px;
			border-radius: 12px;
			font-size: 0.85rem;
			font-weight: 600;
		}
	}

	.description {
		color: var(--muted, #999);
		margin: 8px 0;
		font-size: 0.95rem;
	}

	.preview {
		display: grid;
		grid-template-columns: repeat(10, 1fr);
		gap: 2px;
		margin: 12px 0;
		padding: 8px;
		background: var(--bg);
		border-radius: 6px;

		.preview-cell {
			aspect-ratio: 1;
			display: flex;
			align-items: center;
			justify-content: center;
			background: #2d3748;
			border-radius: 3px;

			.mini-piece {
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				gap: 1px;
				width: 100%;
				height: 100%;
				background: #ff6b6b;
				border-radius: 3px;

				.icon {
					font-size: 0.9rem;
					line-height: 1;
				}

				.rank {
					font-size: 0.5rem;
					font-weight: bold;
					color: white;
					background: rgba(0, 0, 0, 0.3);
					padding: 0px 3px;
					border-radius: 2px;
					line-height: 1;
				}
			}
		}
	}

	.card-actions {
		display: flex;
		gap: 10px;
		margin-top: 16px;

		button {
			flex: 1;
			border: none;
			padding: 8px 16px;
			border-radius: 6px;
			cursor: pointer;
			font-weight: 500;
			transition: background 0.2s;
		}

		.btn-edit {
			background: #4a4;
			color: white;

			&:hover {
				background: #393;
			}
		}

		.btn-delete {
			background: #d44;
			color: white;

			&:hover {
				background: #c33;
			}
		}
	}
</style>
