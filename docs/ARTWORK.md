# Product artwork sources and limitations

Product art is now **backend-controlled**: the Go API returns `products.cover_image_url`. No title-to-asset table or Steam app IDs are embedded in the frontend. The frontend renders only HTTPS URLs from that API field and shows a neutral product title card if absent or unavailable.

`database/004_artwork_metadata.sql` provides title-specific sample image references for the 19 non-Minecraft seeded titles. Each URL follows Steam's public CDN pattern `https://cdn.akamai.steamstatic.com/steam/apps/{app_id}/library_600x900.jpg`. Examples: Portal 2 (620), Hollow Knight (367520), Stardew Valley (413150). `Minecraft: Java Edition` is not mapped and intentionally uses a neutral fallback. See SQL for the complete mapping and provenance.

**Ownership / usage:** Artwork and trademarks belong to the respective game rights holders. Public visibility on the Steam CDN is **not a license** to use those images in an unrelated marketplace or to redistribute them. The URLs are not bundled as image files; the remote host may deny access. Some seeded console-edition entries use title art from a Steam listing, which is **not** a claim that the offer is a real Steam key. Platform, region, stock, seller and price always come from the demo API.

**Before deployment:** verify permitted use for each image, secure appropriate permission/approved assets, or retain neutral title cards. Remove or set `cover_image_url` to NULL for any URL you cannot use. The migration fills blank values only and does not overwrite your existing custom cover images; it never changes user accounts, orders or keys.

**Public demo artwork status (2026-11-11):**
The GameKey internship demo retains the external Steam CDN
cover references for 19 fictional catalogue entries, as defined
in database/004_artwork_metadata.sql.

The sources are documented, but permission or a license to reuse
these images in this marketplace has NOT been verified.
No affiliation, endorsement, or permission from Steam, Eneba,
game publishers, or other rights holders is claimed.

If permission is denied or a rights holder requests removal,
set the relevant cover_image_url values to NULL or replace them
with authorized artwork. The frontend provides neutral title-card
fallbacks. This public-demo decision does not affect the original
local project's database.