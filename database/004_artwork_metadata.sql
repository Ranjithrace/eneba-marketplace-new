-- Adds product cover URLs to the EXISTING API data source. No reseeding or deletions.
-- NOTICE: the URLs are a demonstration reference to exact-title Steam art, NOT a license.
-- Only blank cover_image_url values are filled. Existing custom images remain intact.
-- Review docs/ARTWORK.md before public deployment. Does not modify accounts, orders or keys.
BEGIN;
CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
WITH sources(slug, source_url) AS (VALUES
    ('stardew-valley-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/413150/library_600x900.jpg'),
    ('hades-pc-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/1145360/library_600x900.jpg'),
    ('hollow-knight-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/367520/library_600x900.jpg'),
    ('cyberpunk-2077-pc-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/1091500/library_600x900.jpg'),
    ('the-witcher-3-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/292030/library_600x900.jpg'),
    ('elden-ring-pc-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/1245620/library_600x900.jpg'),
    ('terraria-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/105600/library_600x900.jpg'),
    ('celeste-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/504230/library_600x900.jpg'),
    ('dead-cells-pc-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/588650/library_600x900.jpg'),
    ('subnautica-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/264710/library_600x900.jpg'),
    ('portal-2-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/620/library_600x900.jpg'),
    ('hades-ii-pc-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/1145350/library_600x900.jpg'),
    ('no-mans-sky-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/275850/library_600x900.jpg'),
    ('overcooked-2-switch-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/728880/library_600x900.jpg'),
    ('cuphead-xbox-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/268910/library_600x900.jpg'),
    ('it-takes-two-ps5-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/1426210/library_600x900.jpg'),
    ('untitled-goose-game-switch-global', 'https://cdn.akamai.steamstatic.com/steam/apps/837470/library_600x900.jpg'),
    ('outer-wilds-pc-global', 'https://cdn.akamai.steamstatic.com/steam/apps/753640/library_600x900.jpg'),
    ('slay-the-spire-pc-europe', 'https://cdn.akamai.steamstatic.com/steam/apps/646570/library_600x900.jpg')
)
UPDATE products AS p
SET cover_image_url = sources.source_url
FROM sources
WHERE p.slug = sources.slug AND NULLIF(BTRIM(p.cover_image_url), '') IS NULL;
INSERT INTO schema_migrations(version) VALUES (4) ON CONFLICT DO NOTHING;
COMMIT;
