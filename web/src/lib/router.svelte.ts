let path = $state(window.location.pathname);

window.addEventListener('popstate', () => { path = window.location.pathname; });

export function getPath() { return path; }

export function navigate(to: string) {
    window.history.pushState({}, '', to);
    path = to;
}

export function matchPath(pattern: string): Record<string, string> | null {
    const patParts = pattern.split('/').filter(Boolean);
    const urlParts = path.split('/').filter(Boolean);
    if (patParts.length !== urlParts.length) return null;
    const params: Record<string, string> = {};
    for (let i = 0; i < patParts.length; i++) {
        if (patParts[i].startsWith(':')) {
            params[patParts[i].slice(1)] = urlParts[i];
        } else if (patParts[i] !== urlParts[i]) {
            return null;
        }
    }
    return params;
}
