package catalog

import (
    "strings"
    "testing"
)

// TestBaseQueryNamesCoverImage guards the CTE projection used by /api/products.
// The real PostgreSQL journey is additionally covered by the integration test
// when TEST_DATABASE_URL is configured.
func TestBaseQueryNamesCoverImage(t *testing.T) {
    const projection = "COALESCE(p.cover_image_url,'') AS cover_image_url"
    if !strings.Contains(baseQuery, projection) {
        t.Fatalf("catalog base query must alias cover_image_url for outer CTE selectors")
    }
}
