package queries

const (
	WithSelect = `WITH
		selected_perfumes_base_info AS (
		%s
		),`
	SelectPerfumesBaseInfo = `SELECT
		pb.canonized_brand as canonized_brand,
		pb.canonized_name as canonized_name,
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
		SELECT canonized_brand, canonized_name, sex_id, jsonb_agg(DISTINCT family) FILTER (WHERE family IS NOT NULL) as families
		FROM families
		WHERE (canonized_brand, canonized_name, sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM selected_perfumes_base_info)
		GROUP BY canonized_brand, canonized_name, sex_id
	),
	notes_tags AS (
        SELECT note_name, jsonb_agg(tag_name) as tags
        FROM notes_with_tags
        GROUP BY note_name
    ),
	notes_characteristics AS (
        SELECT note_name, jsonb_agg(jsonb_build_object('name', characteristic_name, 'value', value)) as characteristics
        FROM notes_with_characteristics
        GROUP BY note_name
    ),
    enriched_notes AS (
        SELECT
            COALESCE(nc.note_name, nt.note_name) as note,
            COALESCE(nc.characteristics, '[]'::jsonb) as characteristics,
            COALESCE(nt.tags, '[]'::jsonb) as tags
        FROM notes_characteristics nc
        FULL OUTER JOIN notes_tags nt ON nc.note_name = nt.note_name
    ),
	aggregated_upper_notes AS (
		SELECT 
            un.canonized_brand as canonized_brand,
            un.canonized_name as canonized_name,
            un.sex_id as sex_id,
            jsonb_agg(DISTINCT un.note) FILTER (WHERE un.note IS NOT NULL) as upper_notes,
            jsonb_agg(jsonb_build_object('name', en.note, 'characteristics', COALESCE(en.characteristics, '[]'::jsonb), 'tags', COALESCE(en.tags, '[]'::jsonb))) FILTER (WHERE en.note IS NOT NULL) as enriched_upper_notes
		FROM upper_notes un
        LEFT JOIN enriched_notes en ON un.note = en.note
		WHERE (un.canonized_brand, un.canonized_name, un.sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM selected_perfumes_base_info)
		GROUP BY un.canonized_brand, un.canonized_name, un.sex_id
	),
	aggregated_core_notes AS (
		SELECT
            cn.canonized_brand as canonized_brand,
            cn.canonized_name as canonized_name,
            cn.sex_id as sex_id,
            jsonb_agg(DISTINCT cn.note) FILTER (WHERE cn.note IS NOT NULL) as core_notes,
            jsonb_agg(jsonb_build_object('name', en.note, 'characteristics', COALESCE(en.characteristics, '[]'::jsonb), 'tags', COALESCE(en.tags, '[]'::jsonb))) FILTER (WHERE en.note IS NOT NULL) as enriched_core_notes
		FROM core_notes cn
        LEFT JOIN enriched_notes en ON cn.note = en.note
		WHERE (cn.canonized_brand, cn.canonized_name, cn.sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM selected_perfumes_base_info)
		GROUP BY cn.canonized_brand, cn.canonized_name, cn.sex_id
	),
	aggregated_base_notes AS (
		SELECT 
            bn.canonized_brand as canonized_brand,
            bn.canonized_name as canonized_name,
            bn.sex_id as sex_id,
            jsonb_agg(DISTINCT bn.note) FILTER (WHERE bn.note IS NOT NULL) as base_notes,
            jsonb_agg(jsonb_build_object('name', en.note, 'characteristics', COALESCE(en.characteristics, '[]'::jsonb), 'tags', COALESCE(en.tags, '[]'::jsonb))) FILTER (WHERE en.note IS NOT NULL) as enriched_base_notes
		FROM base_notes bn
        LEFT JOIN enriched_notes en ON bn.note = en.note
		WHERE (bn.canonized_brand, bn.canonized_name, bn.sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM selected_perfumes_base_info)
		GROUP BY bn.canonized_brand, bn.canonized_name, bn.sex_id
	),
	shops_with_variants AS (
		SELECT 
			v.canonized_brand,
			v.canonized_name,
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
		WHERE (v.canonized_brand, v.canonized_name, v.sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM selected_perfumes_base_info)
		GROUP BY v.canonized_brand, v.canonized_name, v.sex_id, sh.id, sh.name, sh.domain
	),
	enriched_selected_perfumes_with_properties AS (
		SELECT
			pb.canonized_brand as canonized_brand,
			pb.canonized_name as canonized_name,
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
				'base_notes', COALESCE(abn.base_notes, '[]'::jsonb),
				'enriched_upper_notes', COALESCE(aun.enriched_upper_notes, '[]'::jsonb),
                'enriched_core_notes', COALESCE(acn.enriched_core_notes, '[]'::jsonb),
                'enriched_base_notes', COALESCE(abn.enriched_base_notes, '[]'::jsonb)
			)::json AS properties
		FROM selected_perfumes_base_info pb
		LEFT JOIN aggregated_families af ON pb.canonized_brand = af.canonized_brand AND pb.canonized_name = af.canonized_name AND pb.sex_id = af.sex_id
		LEFT JOIN aggregated_upper_notes aun ON pb.canonized_brand = aun.canonized_brand AND pb.canonized_name = aun.canonized_name AND pb.sex_id = aun.sex_id
		LEFT JOIN aggregated_core_notes acn ON pb.canonized_brand = acn.canonized_brand AND pb.canonized_name = acn.canonized_name AND pb.sex_id = acn.sex_id
		LEFT JOIN aggregated_base_notes abn ON pb.canonized_brand = abn.canonized_brand AND pb.canonized_name = abn.canonized_name AND pb.sex_id = abn.sex_id
	),
	enriched_shops_with_variants AS (
		SELECT
			swv.canonized_brand as canonized_brand,
			swv.canonized_name as canonized_name,
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
		GROUP BY swv.canonized_brand, swv.canonized_name, swv.sex_id
	)
	SELECT 
		espv.brand as brand,
		espv.name as name,
		espv.sex as sex,
		espv.image_url as image_url,
		espv.properties as properties,
		COALESCE(eswv.shops::json, '[]'::json) as shops
	FROM enriched_selected_perfumes_with_properties espv
	LEFT JOIN enriched_shops_with_variants eswv ON espv.canonized_brand = eswv.canonized_brand AND espv.canonized_name = eswv.canonized_name AND espv.sex_id = eswv.sex_id
	`
)
