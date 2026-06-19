// Inline SVG icons — replaces lucide-react with zero dependencies
type P = { size?: number; className?: string };

const svg = (path: string, p: P) => (
  <svg width={p.size ?? 16} height={p.size ?? 16} viewBox="0 0 24 24"
    fill="none" stroke="currentColor" strokeWidth="2"
    strokeLinecap="round" strokeLinejoin="round" className={p.className}>
    <path d={path} />
  </svg>
);

export const IconPlus         = (p: P) => svg('M12 5v14M5 12h14', p);
export const IconTrash        = (p: P) => svg('M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6', p);
export const IconChevronLeft  = (p: P) => svg('M15 18l-6-6 6-6', p);
export const IconChevronRight = (p: P) => svg('M9 18l6-6-6-6', p);
export const IconChevronDown  = (p: P) => svg('M6 9l6 6 6-6', p);
export const IconChevronUp    = (p: P) => svg('M18 15l-6-6-6 6', p);
export const IconShield       = (p: P) => svg('M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z', p);
export const IconSword        = (p: P) => svg('M14.5 17.5L3 6V3h3l11.5 11.5M13 19l6-6M2 22l2-2M17 3l4 4-9.5 9.5-4-4z', p);
export const IconBook         = (p: P) => svg('M4 19.5A2.5 2.5 0 0 1 6.5 17H20M4 19.5A2.5 2.5 0 0 0 6.5 22H20V2H6.5A2.5 2.5 0 0 0 4 4.5v15z', p);
export const IconBackpack     = (p: P) => svg('M4 20V10a4 4 0 0 1 4-4h8a4 4 0 0 1 4 4v10M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2M8 20h8', p);
export const IconStar         = (p: P) => svg('M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z', p);
export const IconMoon         = (p: P) => svg('M21 12.79A9 9 0 1 1 11.21 3a7 7 0 0 0 9.79 9.79z', p);
export const IconSun          = (p: P) => svg('M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42M12 17a5 5 0 1 0 0-10 5 5 0 0 0 0 10z', p);
export const IconUsers        = (p: P) => svg('M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8zM23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75', p);
export const IconLogOut       = (p: P) => svg('M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9', p);
export const IconSkull        = (p: P) => svg('M12 2a9 9 0 0 1 9 9c0 3.18-1.65 5.97-4.13 7.6L17 21H7l.13-2.4A9 9 0 0 1 3 11a9 9 0 0 1 9-9zM9 17v2M15 17v2M9 12a1 1 0 1 0 0-2 1 1 0 0 0 0 2zM15 12a1 1 0 1 0 0-2 1 1 0 0 0 0 2z', p);
export const IconCheck        = (p: P) => svg('M20 6L9 17l-5-5', p);
export const IconPlay         = (p: P) => svg('M5 3l14 9-14 9V3z', p);
export const IconStop         = (p: P) => svg('M18 6H6v12h12V6z', p);
export const IconSkip         = (p: P) => svg('M5 4l10 8-10 8V4zM19 5v14', p);
