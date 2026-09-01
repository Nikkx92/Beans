package coffee

import (
	"voicer/internal/storage"

	_ "embed"

	"gioui.org/widget"
)

type SettingItem struct {
	Title      string
	Parameters []Parameter
}

type Parameter struct {
	Title  string
	Editor *widget.Editor
}

func NewSettingItems() (
	map[string]map[string]*storage.Type,
	map[string]map[string]*storage.Type) {

	m := make(map[string]map[string]*storage.Type)
	mMobile := make(map[string]map[string]*storage.Type)
	// Проходим по всем зарегистрированным элементам интерфейса
	for _, field := range storage.SettingItemsConfig {

		// Получаем список параметров, которые соответствуют группе этого элемента
		paramTitles, exists := storage.ParamGroups[field.Group]
		if !exists {
			continue // Если группа не описана, пропускаем
		}

		// Создаем карты для текущего элемента
		editorsMap := make(map[string]*storage.Type)
		editorsMapM := make(map[string]*storage.Type)

		// Заполняем дефолтными значениями только те параметры, которые РЕАЛЬНО нужны элементу
		for _, titleParam := range paramTitles {
			editorsMap[titleParam] = &storage.Type{
				Num: 16,
				Str: "16",
			}

			editorsMapM[titleParam] = &storage.Type{
				Num: 16,
				Str: "16",
			}
		}

		if field.IsDesktop {
			m[field.Title] = editorsMap
		}

		if field.IsMobile {
			mMobile[field.Title] = editorsMapM
		}
	}

	return m, mMobile
}

func (c *Components) SetFlags() {
	for i := range c.Blocks {
		block := &c.Blocks[i]
		block.flag = c.flags[block.fieldCountry.Text()]
		c.a.Window.Invalidate()
	}
}

func (c *Components) UpdateState() storage.StateCoffee {
	var state storage.StateCoffee
	state.Items = make([]storage.Items, len(c.Blocks))

	for i, block := range c.Blocks {
		state.Items[i].Country = block.fieldCountry.Text()

		if len(block.Sorts) == 0 {
			continue
		}

		state.Items[i].Sorts = make([]storage.Fields, len(block.Sorts))

		for j, sort := range block.Sorts {
			st := &state.Items[i].Sorts[j]

			st.Sort = sort.fieldSort.Text()
			st.Blend = sort.fieldBlend.Text()
			st.Grade = sort.fieldGrade.Text()
			st.Location = sort.fieldLocation.Text()
			st.HeightAbove = sort.fieldHeightAbove.Text()
			st.Drying = sort.fieldDrying.Text()
			st.CornSize = sort.fieldCornSize.Text()
			st.Taste = sort.fieldTaste.Text()
			st.Price = sort.fieldPrice.Text()
			st.Quantity = sort.fieldQuantity.Text()
			st.StockDesc = sort.fieldStockDesc.Text()
			st.InStock = sort.inStock.Value
			st.Limited = sort.limited.Value
		}

	}

	state.Parameters = updParameters(c.Settings.SettingItems, c.Settings.FontSz)
	state.MobileParameters = updParameters(c.Settings.SettingItems, c.Settings.FontSzM)

	c.a.State.StateCoffee = state
	return state
}

func updParameters(items []SettingItem, m map[string]map[string]*storage.Type) []storage.FieldParameters {
	params := make([]storage.FieldParameters, len(items))
	for i := range params {
		s := &params[i]
		s.Title = items[i].Title

		if subMap, ok := m[s.Title]; ok {
			s.Parameters = make(map[string]storage.Type)
			for param, val := range subMap {
				s.Parameters[param] = *val
			}
		}
	}
	return params
}

func (c *Components) SetState(s storage.StateCoffee) {
	for i, state := range s.Items {
		c.Blocks[i] = addItemBlock()
		c.Blocks[i].fieldCountry.SetText(state.Country)

		if len(state.Sorts) == 0 {
			continue
		}

		c.Blocks[i].Sorts = make([]SortBlock, len(state.Sorts))

		for j, sortState := range state.Sorts {
			c.Blocks[i].Sorts[j] = addSortBlock()
			sb := &c.Blocks[i].Sorts[j]

			sb.fieldSort.SetText(sortState.Sort)
			sb.fieldBlend.SetText(sortState.Blend)
			sb.fieldGrade.SetText(sortState.Grade)
			sb.fieldLocation.SetText(sortState.Location)
			sb.fieldHeightAbove.SetText(sortState.HeightAbove)
			sb.fieldDrying.SetText(sortState.Drying)
			sb.fieldCornSize.SetText(sortState.CornSize)
			sb.fieldTaste.SetText(sortState.Taste)
			sb.fieldPrice.SetText(sortState.Price)
			sb.fieldQuantity.SetText(sortState.Quantity)
			sb.fieldStockDesc.SetText(sortState.StockDesc)
			sb.inStock.Value = sortState.InStock
			sb.limited.Value = sortState.Limited
		}
	}

	setParameters(s.Parameters, c.Settings.FontSz)
	setParameters(s.MobileParameters, c.Settings.FontSzM)

}

func setParameters(p []storage.FieldParameters, m map[string]map[string]*storage.Type) {
	for _, v := range p {
		m[v.Title] = make(map[string]*storage.Type)
		for param, val := range v.Parameters {
			m[v.Title][param] = &storage.Type{
				Num: val.Num,
				Str: val.Str,
			}
		}
	}
}
