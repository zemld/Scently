package queries

const (
	With = `WITH
		selected_perfumes_base_info AS (
		%s
		),`
	PerfumesBaseInfo = `SELECT
		pb.brand as brand,
		pb.name as name,
		pb.sex_id as sex_id,
		s.sex as sex,
		pb.type as type,
		pb.image_url as image_url
	FROM perfume_base_info pb
	INNER JOIN sexes s ON pb.sex_id = s.id
	`

	EnrichSelectedPerfumes = `
	aggregated_families AS (
		SELECT brand, name, sex_id, jsonb_agg(DISTINCT family) FILTER (WHERE family IS NOT NULL) as families
		FROM families
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM selected_perfumes_base_info)
		GROUP BY brand, name, sex_id
	),
	aggregated_upper_notes AS (
		SELECT brand, name, sex_id, jsonb_agg(DISTINCT note) FILTER (WHERE note IS NOT NULL) as upper_notes
		FROM upper_notes
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM selected_perfumes_base_info)
		GROUP BY brand, name, sex_id
	),
	aggregated_core_notes AS (
		SELECT brand, name, sex_id, jsonb_agg(DISTINCT note) FILTER (WHERE note IS NOT NULL) as core_notes
		FROM core_notes
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM selected_perfumes_base_info)
		GROUP BY brand, name, sex_id
	),
	aggregated_base_notes AS (
		SELECT brand, name, sex_id, jsonb_agg(DISTINCT note) FILTER (WHERE note IS NOT NULL) as base_notes
		FROM base_notes
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM selected_perfumes_base_info)
		GROUP BY brand, name, sex_id
	),
	shops_with_variants AS (
		SELECT 
			v.brand,
			v.name,
			v.sex_id,
			sh.id as shop_id,
			sh.name as shop_name,
			sh.domain as domain,
			jsonb_agg(
				jsonb_build_object(
					'volume', v.volume,
					'price', v.price,
					'link', v.link
				)
			) FILTER (WHERE v.volume IS NOT NULL) as variants
		FROM variants v
		INNER JOIN shops sh ON v.shop_id = sh.id
		WHERE (v.brand, v.name, v.sex_id) IN (SELECT brand, name, sex_id FROM selected_perfumes_base_info)
		GROUP BY v.brand, v.name, v.sex_id, sh.id, sh.name, sh.domain
	),
	enriched_selected_perfumes_with_properties AS (
		SELECT
			pb.brand as brand,
			pb.name as name,
			pb.sex_id as sex_id,
			pb.sex as sex,
			pb.image_url as image_url,
			jsonb_build_object(
				'perfume_type', pb.type,
				'family', COALESCE(af.families, '[]'::jsonb),
				'upper_notes', COALESCE(aun.upper_notes, '[]'::jsonb),
				'core_notes', COALESCE(acn.core_notes, '[]'::jsonb),
				'base_notes', COALESCE(abn.base_notes, '[]'::jsonb)
			)::json AS properties
		FROM selected_perfumes_base_info pb
		LEFT JOIN aggregated_families af ON pb.brand = af.brand AND pb.name = af.name AND pb.sex_id = af.sex_id
		LEFT JOIN aggregated_upper_notes aun ON pb.brand = aun.brand AND pb.name = aun.name AND pb.sex_id = aun.sex_id
		LEFT JOIN aggregated_core_notes acn ON pb.brand = acn.brand AND pb.name = acn.name AND pb.sex_id = acn.sex_id
		LEFT JOIN aggregated_base_notes abn ON pb.brand = abn.brand AND pb.name = abn.name AND pb.sex_id = abn.sex_id
	),
	enriched_shops_with_variants AS (
		SELECT
			swv.brand as brand,
			swv.name as name,
			swv.sex_id as sex_id,
			jsonb_agg(
				jsonb_build_object(
					'shop_name', swv.shop_name,
					'domain', swv.domain,
					'image_url', NULL,
					'variants', swv.variants
				)
			) FILTER (WHERE swv.shop_id IS NOT NULL) AS shops
		FROM shops_with_variants swv
		GROUP BY swv.brand, swv.name, swv.sex_id
	)
	SELECT 
		espv.brand as brand,
		espv.name as name,
		espv.sex as sex,
		espv.image_url as image_url,
		espv.properties as properties,
		COALESCE(eswv.shops::json, '[]'::json) as shops
	FROM enriched_selected_perfumes_with_properties espv
	LEFT JOIN enriched_shops_with_variants eswv ON espv.brand = eswv.brand AND espv.name = eswv.name AND espv.sex_id = eswv.sex_id
	`
)
