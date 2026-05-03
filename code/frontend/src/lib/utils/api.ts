/**
 * Parses an API-formatted message.
 * Supports both JSON strings and objects with keys like 'error', 'warning', etc.
 */
export function parseApiMessage(data: any): { message: string; type: string } {
    if (!data) return { message: "", type: "error" };

    // If it's already an object (e.g. from a websocket or already parsed JSON)
    if (typeof data === 'object' && !Array.isArray(data)) {
        const key = Object.keys(data)[0];
        const message = data[key];
        if (typeof message === 'string') {
            return { message, type: key };
        }
        
        // If the object doesn't match our expected format, just stringify it
        // but avoid returning empty JSON objects as messages
        const stringified = JSON.stringify(data);
        return { 
            message: stringified === '{}' ? "" : stringified, 
            type: 'error' 
        };
    }

    // If it's a string, try to parse it as JSON
    if (typeof data === 'string') {
        try {
            const json = JSON.parse(data);
            return parseApiMessage(json);
        } catch {
            // Not JSON, use as plain text
            return { message: data, type: "error" };
        }
    }

    return { message: String(data), type: "error" };
}