-- Stage 2: 20 sample products, fictional sellers and demonstration offers.
-- Re-running will not duplicate products, sellers or offers.
BEGIN;
INSERT INTO sellers(name) VALUES ('North Star Games'), ('Pixel Harbor'), ('Key Lantern') ON CONFLICT (name) DO NOTHING;

INSERT INTO products(slug,title,description,platform,region,genre) VALUES
    ('minecraft-java-pc-global', 'Minecraft: Java Edition', 'Sandbox building and survival game.', 'PC', 'Global', 'Sandbox'),
    ('stardew-valley-pc-global', 'Stardew Valley', 'Relaxing farming and community adventure.', 'PC', 'Global', 'Simulation'),
    ('hades-pc-europe', 'Hades', 'Fast-paced roguelike action.', 'PC', 'Europe', 'Action'),
    ('hollow-knight-pc-global', 'Hollow Knight', 'Atmospheric exploration platformer.', 'PC', 'Global', 'Adventure'),
    ('cyberpunk-2077-pc-europe', 'Cyberpunk 2077', 'Open-world science-fiction role-playing game.', 'PC', 'Europe', 'RPG'),
    ('the-witcher-3-pc-global', 'The Witcher 3: Wild Hunt', 'Story-driven fantasy RPG.', 'PC', 'Global', 'RPG'),
    ('elden-ring-pc-europe', 'Elden Ring', 'Open-world fantasy action role-playing game.', 'PC', 'Europe', 'RPG'),
    ('terraria-pc-global', 'Terraria', '2D sandbox exploration and crafting.', 'PC', 'Global', 'Sandbox'),
    ('celeste-pc-global', 'Celeste', 'Precision platforming adventure.', 'PC', 'Global', 'Platformer'),
    ('dead-cells-pc-europe', 'Dead Cells', 'Action-platformer with roguelike elements.', 'PC', 'Europe', 'Action'),
    ('subnautica-pc-global', 'Subnautica', 'Underwater survival and exploration.', 'PC', 'Global', 'Survival'),
    ('portal-2-pc-global', 'Portal 2', 'First-person puzzle adventure.', 'PC', 'Global', 'Puzzle'),
    ('hades-ii-pc-europe', 'Hades II', 'Roguelike action adventure.', 'PC', 'Europe', 'Action'),
    ('no-mans-sky-pc-global', 'No Man''s Sky', 'Space exploration and survival.', 'PC', 'Global', 'Adventure'),
    ('overcooked-2-switch-europe', 'Overcooked! 2', 'Cooperative cooking game.', 'Nintendo Switch', 'Europe', 'Party'),
    ('cuphead-xbox-europe', 'Cuphead', 'Run-and-gun action with a vintage cartoon aesthetic.', 'Xbox', 'Europe', 'Action'),
    ('it-takes-two-ps5-europe', 'It Takes Two', 'Cooperative story-driven platform adventure.', 'PlayStation 5', 'Europe', 'Adventure'),
    ('untitled-goose-game-switch-global', 'Untitled Goose Game', 'Playful stealth puzzle sandbox.', 'Nintendo Switch', 'Global', 'Puzzle'),
    ('outer-wilds-pc-global', 'Outer Wilds', 'Exploration mystery in a time-loop solar system.', 'PC', 'Global', 'Adventure'),
    ('slay-the-spire-pc-europe', 'Slay the Spire', 'Deckbuilding roguelike.', 'PC', 'Europe', 'Strategy')
ON CONFLICT (slug) DO NOTHING;

-- Up to three separate sellers can offer the same product.
WITH price_data(slug, seller_name, price_cents) AS (
    VALUES
        ('minecraft-java-pc-global', 'North Star Games', 799),
        ('minecraft-java-pc-global', 'Pixel Harbor', 1029),
        ('minecraft-java-pc-global', 'Key Lantern', 714),
        ('stardew-valley-pc-global', 'North Star Games', 909),
        ('stardew-valley-pc-global', 'Pixel Harbor', 1139),
        ('stardew-valley-pc-global', 'Key Lantern', 824),
        ('hades-pc-europe', 'North Star Games', 1019),
        ('hades-pc-europe', 'Pixel Harbor', 1249),
        ('hades-pc-europe', 'Key Lantern', 934),
        ('hollow-knight-pc-global', 'North Star Games', 1129),
        ('hollow-knight-pc-global', 'Pixel Harbor', 1359),
        ('hollow-knight-pc-global', 'Key Lantern', 1044),
        ('cyberpunk-2077-pc-europe', 'North Star Games', 1239),
        ('cyberpunk-2077-pc-europe', 'Pixel Harbor', 1469),
        ('cyberpunk-2077-pc-europe', 'Key Lantern', 1154),
        ('the-witcher-3-pc-global', 'North Star Games', 1349),
        ('the-witcher-3-pc-global', 'Pixel Harbor', 1579),
        ('the-witcher-3-pc-global', 'Key Lantern', 1264),
        ('elden-ring-pc-europe', 'North Star Games', 1459),
        ('elden-ring-pc-europe', 'Pixel Harbor', 1689),
        ('elden-ring-pc-europe', 'Key Lantern', 1374),
        ('terraria-pc-global', 'North Star Games', 1569),
        ('terraria-pc-global', 'Pixel Harbor', 1799),
        ('terraria-pc-global', 'Key Lantern', 1484),
        ('celeste-pc-global', 'North Star Games', 1679),
        ('celeste-pc-global', 'Pixel Harbor', 1909),
        ('celeste-pc-global', 'Key Lantern', 1594),
        ('dead-cells-pc-europe', 'North Star Games', 1789),
        ('dead-cells-pc-europe', 'Pixel Harbor', 2019),
        ('dead-cells-pc-europe', 'Key Lantern', 1704),
        ('subnautica-pc-global', 'North Star Games', 1899),
        ('subnautica-pc-global', 'Pixel Harbor', 2129),
        ('subnautica-pc-global', 'Key Lantern', 1814),
        ('portal-2-pc-global', 'North Star Games', 2009),
        ('portal-2-pc-global', 'Pixel Harbor', 2239),
        ('portal-2-pc-global', 'Key Lantern', 1924),
        ('hades-ii-pc-europe', 'North Star Games', 2119),
        ('hades-ii-pc-europe', 'Pixel Harbor', 2349),
        ('hades-ii-pc-europe', 'Key Lantern', 2034),
        ('no-mans-sky-pc-global', 'North Star Games', 2229),
        ('no-mans-sky-pc-global', 'Pixel Harbor', 2459),
        ('no-mans-sky-pc-global', 'Key Lantern', 2144),
        ('overcooked-2-switch-europe', 'North Star Games', 2339),
        ('overcooked-2-switch-europe', 'Pixel Harbor', 2569),
        ('overcooked-2-switch-europe', 'Key Lantern', 2254),
        ('cuphead-xbox-europe', 'North Star Games', 2449),
        ('cuphead-xbox-europe', 'Pixel Harbor', 2679),
        ('cuphead-xbox-europe', 'Key Lantern', 2364),
        ('it-takes-two-ps5-europe', 'North Star Games', 2559),
        ('it-takes-two-ps5-europe', 'Pixel Harbor', 2789),
        ('it-takes-two-ps5-europe', 'Key Lantern', 2474),
        ('untitled-goose-game-switch-global', 'North Star Games', 2669),
        ('untitled-goose-game-switch-global', 'Pixel Harbor', 2899),
        ('untitled-goose-game-switch-global', 'Key Lantern', 2584),
        ('outer-wilds-pc-global', 'North Star Games', 2779),
        ('outer-wilds-pc-global', 'Pixel Harbor', 3009),
        ('outer-wilds-pc-global', 'Key Lantern', 2694),
        ('slay-the-spire-pc-europe', 'North Star Games', 2889),
        ('slay-the-spire-pc-europe', 'Pixel Harbor', 3119),
        ('slay-the-spire-pc-europe', 'Key Lantern', 2804)
)
INSERT INTO offers(product_id, seller_id, price_cents)
SELECT p.id, s.id, pd.price_cents
FROM price_data pd
JOIN products p ON p.slug = pd.slug
JOIN sellers s ON s.name = pd.seller_name
ON CONFLICT (product_id, seller_id) DO NOTHING;

INSERT INTO homepage_sections(section_key, heading, subtitle, display_order) VALUES
    ('banner', 'Discover your next game', 'Browse our demo digital-key marketplace', 1),
    ('bestsellers', 'Best sellers', 'Popular picks in our demo catalog', 2),
    ('deals', 'Featured deals', 'Compare sample seller prices', 3)
ON CONFLICT (section_key) DO NOTHING;

-- Two non-redeemable keys per offer, reserved for the later checkout stage.
INSERT INTO digital_keys (offer_id, key_code)
SELECT o.id, 'DEMO-ONLY-' || LPAD(o.id::TEXT, 5, '0') || '-' || LPAD(n.seq::TEXT, 2, '0')
FROM offers o
CROSS JOIN generate_series(1, 2) AS n(seq)
ON CONFLICT (key_code) DO NOTHING;

COMMIT;
