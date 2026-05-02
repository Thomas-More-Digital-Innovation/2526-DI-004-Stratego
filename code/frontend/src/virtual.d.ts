declare module 'virtual:changelog' {
	const content: string;
	export default content;
}

// Manual Node.js type declarations to fix compiler errors when @types/node is missing
declare module 'node:fs' {
	export function readFileSync(path: string, options?: any): any;
}

declare module 'node:url' {
	export function fileURLToPath(url: URL | string): string;
}

declare module 'node:path' {
	export function resolve(...paths: string[]): string;
}
