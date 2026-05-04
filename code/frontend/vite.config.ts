import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, type Plugin } from 'vite';
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const pkg = JSON.parse(readFileSync(fileURLToPath(new URL('package.json', import.meta.url)), 'utf8'));

// Virtual module plugin to load the changelog content at build time
const changelogPlugin = (): Plugin => {
    const virtualModuleId = 'virtual:changelog';
    const resolvedVirtualModuleId = '\0' + virtualModuleId;

    return {
        name: 'changelog-loader',
        resolveId(id: string) {
            if (id === virtualModuleId) {
                return resolvedVirtualModuleId;
            }
            return null;
        },
        load(id: string) {
            if (id === resolvedVirtualModuleId) {
                const configDir = fileURLToPath(new URL('.', import.meta.url));
                const changelogPath = path.resolve(configDir, '../../CHANGELOG.md');
                const content = readFileSync(changelogPath, 'utf8');
                return `export default ${JSON.stringify(content)};`;
            }
            return null;
        }
    };
};

export default defineConfig({
    plugins: [tailwindcss(), sveltekit(), changelogPlugin()],
    define: {
        __VERSION__: JSON.stringify(pkg.version)
    }
});
