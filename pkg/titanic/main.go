package titanic

import (
	"context"
	"fmt"
	imports "github.com/rocketlaunchr/dataframe-go/imports"
	chart "github.com/wcharczuk/go-chart"
	plot "github.com/rocketlaunchr/dataframe-go/plot"
	dataframe "github.com/rocketlaunchr/dataframe-go"
	"strconv"
	"sort"
	"os"
)

var (
	sibsp int
	parch int
	sibspparch string
)

func Run(){

	ctx := context.Background()

	csv, err := os.Open("assets/titanic.csv")
	if err != nil {
		Throw(err, "Error opening dataset file")
	}
	defer csv.Close()

	dataset, err := imports.LoadFromCSV(ctx, csv)
	if err != nil {
		Throw(err, "Error load dataset file")
	}
	//============================= PRINT ALL DATASET ===============================
	fmt.Print(dataset.Table())

	//============================= PRINT ROW COUNT =================================
	fmt.Println("Jumlah baris:", dataset.NRows())

	//============================= PRINT COLLUMN COUNT =============================
	fmt.Println("Jumlah kolom:", len(dataset.Series))

	//============================= PRINT SPECIFIC DATASET ==========================
	data := dataframe.NewDataFrame(
		dataset.Series[dataset.MustNameToColumn("Name")],
		dataset.Series[dataset.MustNameToColumn("Sex")],
		dataset.Series[dataset.MustNameToColumn("Age")],
		dataset.Series[dataset.MustNameToColumn("Pclass")],
		dataset.Series[dataset.MustNameToColumn("Fare")],
	)
	fmt.Print(data.Table())

	//============================= PRINT SURVIVAL DATASET ===========================
	class := dataframe.NewDataFrame(
		dataset.Series[dataset.MustNameToColumn("Survived")],
	)
	fmt.Print(class.Table())

	//============================= ADD DATASET WITH Relative FEATURE ================
	relatives := make([]string, dataset.NRows())
	iterator := dataset.ValuesIterator(dataframe.ValuesOptions{
		InitialRow:0,
		Step:1,
		DontReadLock:true,
	})
	dataset.Lock()
	for {
		row, val, _ := iterator()
		if row == nil {
			break
		}
		sibsp,_= strconv.Atoi(val[6].(string))
		parch,_= strconv.Atoi(val[7].(string))
		sibspparch = strconv.Itoa(sibsp+parch)
		relatives[*row] = sibspparch
	}
	dataset.Unlock()
	sr := dataframe.NewSeriesString("Relatives", nil, relatives)
	dataset.AddSeries(sr, nil)
	fmt.Print(dataset.Table())

	//============================ FIND PASSANGER by PCLASS =========================
	pclass := map[string]int{}
	pclassIterator := dataset.ValuesIterator(dataframe.ValuesOptions{
		InitialRow:0,
		Step:1,
		DontReadLock:true,
	})
	dataset.Lock()
	for {
		row, val, _ := pclassIterator()
		if row == nil {
			break
		}
		//============= ADD Pclass class ==================
		if len(pclass) == 0 {
			pclass[val[2].(string)] = 0
		} else {
			if !CheckDuplicateValue(pclass, val[2].(string)) {
				pclass[val[2].(string)] = 0
			}
		}
		//============= Count Passangers ==================
		pclass[val[2].(string)] += 1
	}
	dataset.Unlock()
	passangers := 0
	for key, val := range pclass{
		fmt.Printf("Pclass %s:%d\n", key, val)
		passangers += val
	}
	fmt.Println("Total passanger is:", passangers)

	//============================ FIND PASSANGER by SEX =========================
	psex := map[string]int{}
	psexIterator := dataset.ValuesIterator(dataframe.ValuesOptions{
		InitialRow:0,
		Step:1,
		DontReadLock:true,
	})
	dataset.Lock()
	for {
		row, val, _ := psexIterator()
		if row == nil {
			break
		}
		//============= ADD Pclass class ==================
		if len(psex) == 0 {
			psex[val[4].(string)] = 0
		} else {
			if !CheckDuplicateValue(psex, val[4].(string)) {
				psex[val[4].(string)] = 0
			}
		}
		//============= Count Passangers ==================
		psex[val[4].(string)] += 1
	}
	dataset.Unlock()
	passangers = 0
	for key, val := range psex{
		fmt.Printf("Gender %s:%d\n", key, val)
		passangers += val
	}
	fmt.Println("Total passanger is:", passangers)

	//============================ FIND PASSANGER by AGE =========================
	page := map[string]int{}
	pageIterator := dataset.ValuesIterator(dataframe.ValuesOptions{
		InitialRow:0,
		Step:1,
		DontReadLock:true,
	})
	dataset.Lock()
	for {
		row, val, _ := pageIterator()
		if row == nil {
			break
		}
		//============= ADD Page class ==================
		if len(page) == 0 {
			page[val[5].(string)] = 0
		} else {
			if !CheckDuplicateValue(page, val[5].(string)){
				if val[5].(string) != "" {
					page[val[5].(string)] = 0
				}
			}
		}
		//============= Count Passangers ==================
		if val[5].(string) != "" {
			page[val[5].(string)] += 1
		}
	}
	dataset.Unlock()
	passangers = 0
	for key, val := range page{
		fmt.Printf("Age %s:%d\n", key, val)
		passangers += val
	}
	fmt.Println("Total passanger is:", passangers)

	//============================ SURVIVED PASSENGER =========================
	fmt.Println()
	pcl := []string{"1","2","3"}
	for _, v := range pcl{
		filterFn := dataframe.FilterDataFrameFn(func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error){
			if vals["Survived"] == "1" && vals["Pclass"] == v {
				return dataframe.KEEP, nil
			}
			return dataframe.DROP, nil
		})
		survived, _ := dataframe.Filter(ctx, dataset, filterFn, dataframe.FilterOptions{InPlace:false})
		fmt.Printf("Passanger survived in class[%s]: %d people\n", v, survived.(*dataframe.DataFrame).NRows())
		//=========================
	}
	fmt.Println()
	for _, v := range pcl{
		filterFn := dataframe.FilterDataFrameFn(func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error){
			if vals["Survived"] == "0" && vals["Pclass"] == v {
				return dataframe.KEEP, nil
			}
			return dataframe.DROP, nil
		})
		survived, _ := dataframe.Filter(ctx, dataset, filterFn, dataframe.FilterOptions{InPlace:false})
		fmt.Printf("Passanger not survived in class[%s]: %d people\n", v, survived.(*dataframe.DataFrame).NRows())
	}
	fmt.Println()

	PrintDataBySex(psex)
	PrintDataByAge(page)
}

func PrintDataBySex(psex map[string]int){
	barsValue := []chart.Value{}
	for k, v := range psex{
		vString := strconv.Itoa(v)
		vFloat,_ := strconv.ParseFloat(vString, 64)
		barsValue = append(barsValue, chart.Value{
			Value: vFloat,
			Label: k+" "+vString,
		})
	}

	graph := chart.BarChart{
		Title: "Titanic Passanger Data",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		BarWidth: 200,
		Bars: barsValue,
		YAxis: chart.YAxis{
			Name: "Jumlah",
			Style: chart.Style{
				Hidden: false,
			},
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: 700,
			},
		},
	}

	plt, _ := plot.Open("Titanic Passanger Data", 1366, 768)
	graph.Render(chart.SVG, plt)
	plt.Display()
	<-plt.Closed
}

func PrintDataByAge(page map[string]int){
	barsValue := []chart.Value{}
	i := 0

	var sortedKeys = []float64{}
	for k,_ := range page{
		floatV,_ := strconv.ParseFloat(k, 64)
		sortedKeys = append(sortedKeys, floatV)
	}
	sort.Float64s(sortedKeys)

	for _, v := range sortedKeys{
		vString := strconv.Itoa(page[fmt.Sprintf("%g", v)])
		vFloat,_ := strconv.ParseFloat(vString, 64)
		barsValue = append(barsValue, chart.Value{
			Value: vFloat,
			Label: fmt.Sprintf("%g", v),
		})
		i++
		if i == 20 {
			graph := chart.BarChart{
				Title: "Titanic Passanger Data by AGE",
				Background: chart.Style{
					Padding: chart.Box{
						Top: 40,
					},
				},
				Bars: barsValue,
				YAxis: chart.YAxis{
					Range: &chart.ContinuousRange{
						Min: 0,
						Max: 100,
					},
				},
			}

			plt, _ := plot.Open("Titanic Passanger Data", 1366, 768)
			graph.Render(chart.SVG, plt)
			plt.Display()
			<-plt.Closed

			barsValue = []chart.Value{}
			i = 0
		}
	}
}
