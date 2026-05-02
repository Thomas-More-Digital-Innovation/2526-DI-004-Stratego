import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import path from 'path';

const pkg = JSON.parse(readFileSync(fileURLToPath(new URL('package.json', import.meta.url)), 'utf8'));

// Virtual module plugin to load the changelog content at build time
const changelogPlugin = () => {
    const virtualModuleId = 'virtual:changelog';
    const resolvedVirtualModuleId = '\0' + virtualModuleId;

    return {
        name: 'changelog-loader',
        resolveId(id) {
            if (id === virtualModuleId) {
                return resolvedVirtualModuleId;
            }
        },
        load(id) {
            if (id === resolvedVirtualModuleId) {
                const configDir = fileURLToPath(new URL('.', import.meta.url));
                const changelogPath = path.resolve(configDir, './src/lib/components/changelog/CHANGELOG.md');
                const content = readFileSync(changelogPath, 'utf8');
                return `export default ${JSON.stringify(content)};`;
            }
        }
    };
};

export default defineConfig({
    plugins: [tailwindcss(), sveltekit(), changelogPlugin()],
    define: {
        __VERSION__: JSON.stringify(pkg.version)
    }
});
