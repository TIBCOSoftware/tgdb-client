package samples

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"time"
)

// Data attributes about each person
type PersonData struct {
	MemberName string
	CrownName  string
	HouseHead  bool
	YearBorn   int
	YearDied   int
	ReignStart string
	ReignEnd   string
	CrownTitle string
}

// Members of the House to be inserted in the database
var HouseMemberData = []PersonData{
	PersonData{"Carlo Bonaparte", "", false, 1746, 1785, "", "", ""},
	PersonData{"Letizia Ramolino", "", false, 1750, 1836, "", "", ""},
	PersonData{"Joseph Bonaparte", "Joseph I", false, 1768, 1844, "6 Jun 1808", "11 Dec 1813", "King of Spain"},
	PersonData{"Napoleon Bonaparte", "Napoleon I", false, 1769, 1821, "18 May 1804", "22 Jun 1815", "Emperor of the French"},
	PersonData{"Lucien Bonaparte", "", false, 1775, 1840, "", "", ""},
	PersonData{"Elisa Bonaparte", "Elisa Bonaparte", false, 1777, 1820, "3 Mar 1809", "1 Feb 1814", "Grand Duchess of Tuscany"},
	PersonData{"Louis Bonaparte", "Louis I", false, 1778, 1846, "5 Jun 1806", "1 Jul 1810", "King of Holland"},
	PersonData{"Pauline Bonaparte", "", false, 1780, 1825, "", "", ""},
	PersonData{"Caroline Bonaparte", "", false, 1782, 1839, "", "", ""},
	PersonData{"Jerome Bonaparte", "Jerome I", false, 1784, 1860, "8 Jul 1807", "26 Oct 1813", "King of Westphalia"},
	PersonData{"Marie Louise of Austria", "", false, 1791, 1847, "", "", "Empress Consort of the French"},
	PersonData{"Josephine of Beauharnais", "", false, 1763, 1814, "", "", "Empress Consort of the French"},
	PersonData{"Alexandre of Beauharnais", "", false, 1760, 1794, "", "", ""},
	PersonData{"Betsy Patterson", "", false, 1785, 1879, "", "", ""},
	PersonData{"Catharina of Wurttemberg", "", false, 1783, 1835, "", "", "Queen Consort of Westphalia"},
	PersonData{"Francois Bonaparte", "Napoleon II", false, 1811, 1832, "22 Jun 1815", "7 Jul 1815", "Emperor of the French"},
	PersonData{"Hortense of Beauharnais", "", false, 1783, 1837, "", "", "Queen Consort of Holland"},
	PersonData{"Jerome Napoleon", "", false, 1805, 1870, "", "", ""},
	PersonData{"Prince Napoleon", "", false, 1822, 1891, "", "", ""},
	PersonData{"Louis Napoleon", "Napoleon III", true, 1808, 1873, "2 Dec 1852", "4 Sep 1870", "Emperors of the French"},
	PersonData{"Napoleon-Louis Bonaparte", "Louis II", false, 1804, 1831, "1 Jul 1810", "13 Jul 1810", "King of Holland"},
	PersonData{"Napoleon IV Eugene", "", true, 1856, 1879, "", "", ""},
	PersonData{"Napoleon V Victor", "", true, 1862, 1926, "", "", ""},
	PersonData{"Marie Clotilde Bonaparte", "", false, 1912, 1996, "", "", ""},
	PersonData{"Napoleon VI Louis", "", true, 1914, 1997, "", "", ""},
	PersonData{"Napoleon VII Charles", "", true, 1950, 1999, "", "", ""},
	PersonData{"Napoleon VIII Jean-Christophe", "", true, 1986, 2018, "", "", ""},
	PersonData{"Sophie Catherine Bonaparte", "", false, 1992, -1, "", "", ""},
}

// Relation among the members of the House
type PersonRelation struct {
	FromMemberName string
	ToMemberName   string
	RelationDesc   string
	Attribute1     int    // This could be marriage year if relationDesc = 'spouse', else child's order of birth
	Attribute2     string // This could be marriage city if relationDesc = 'spouse', else EMPTY string
}

// Relation among the members of the House
var HouseRelationData = []PersonRelation{
	PersonRelation{"Carlo Bonaparte", "Letizia Ramolino", "spouse", 1529, "Paris"},
	PersonRelation{"Carlo Bonaparte", "Joseph Bonaparte", "child", 1, ""},
	PersonRelation{"Letizia Ramolino", "Joseph Bonaparte", "child", 1, ""},
	PersonRelation{"Carlo Bonaparte", "Napoleon Bonaparte", "child", 2, ""},
	PersonRelation{"Letizia Ramolino", "Napoleon Bonaparte", "child", 2, ""},
	PersonRelation{"Carlo Bonaparte", "Lucien Bonaparte", "child", 3, ""},
	PersonRelation{"Letizia Ramolino", "Lucien Bonaparte", "child", 3, ""},
	PersonRelation{"Carlo Bonaparte", "Elisa Bonaparte", "child", 4, ""},
	PersonRelation{"Letizia Ramolino", "Elisa Bonaparte", "child", 4, ""},
	PersonRelation{"Carlo Bonaparte", "Louis Bonaparte", "child", 5, ""},
	PersonRelation{"Letizia Ramolino", "Louis Bonaparte", "child", 5, ""},
	PersonRelation{"Carlo Bonaparte", "Pauline Bonaparte", "child", 6, ""},
	PersonRelation{"Letizia Ramolino", "Pauline Bonaparte", "child", 6, ""},
	PersonRelation{"Carlo Bonaparte", "Caroline Bonaparte", "child", 7, ""},
	PersonRelation{"Letizia Ramolino", "Caroline Bonaparte", "child", 7, ""},
	PersonRelation{"Carlo Bonaparte", "Jerome Bonaparte", "child", 8, ""},
	PersonRelation{"Letizia Ramolino", "Jerome Bonaparte", "child", 8, ""},
	PersonRelation{"Napoleon Bonaparte", "Marie Louise of Austria", "spouse", 1656, "Lyon"},
	PersonRelation{"Napoleon Bonaparte", "Francois Bonaparte", "child", 1, ""},
	PersonRelation{"Marie Louise of Austria", "Francois Bonaparte", "child", 1, ""},
	PersonRelation{"Napoleon Bonaparte", "Josephine of Beauharnais", "spouse", 1535, "Cannes"},
	PersonRelation{"Alexandre of Beauharnais", "Josephine of Beauharnais", "spouse", 1705, "Paris"},
	PersonRelation{"Alexandre of Beauharnais", "Hortense of Beauharnais", "child", 1, ""},
	PersonRelation{"Josephine of Beauharnais", "Hortense of Beauharnais", "child", 1, ""},
	PersonRelation{"Louis Bonaparte", "Hortense of Beauharnais", "spouse", 1715, "Paris"},
	PersonRelation{"Louis Bonaparte", "Louis Napoleon", "child", 1, ""},
	PersonRelation{"Hortense of Beauharnais", "Louis Napoleon", "child", 1, ""},
	PersonRelation{"Louis Bonaparte", "Napoleon-Louis Bonaparte", "child", 2, ""},
	PersonRelation{"Hortense of Beauharnais", "Napoleon-Louis Bonaparte", "child", 2, ""},
	PersonRelation{"Jerome Bonaparte", "Betsy Patterson", "spouse", 1669, "Lyon"},
	PersonRelation{"Jerome Bonaparte", "Jerome Napoleon", "child", 1, ""},
	PersonRelation{"Betsy Patterson", "Jerome Napoleon", "child", 1, ""},
	PersonRelation{"Jerome Bonaparte", "Catharina of Wurttemberg", "spouse", 1801, "Nice"},
	PersonRelation{"Jerome Bonaparte", "Prince Napoleon", "child", 1, ""},
	PersonRelation{"Catharina of Wurttemberg", "Prince Napoleon", "child", 1, ""},
	PersonRelation{"Louis Napoleon", "Napoleon IV Eugene", "child", 1, ""},
	PersonRelation{"Prince Napoleon", "Napoleon V Victor", "child", 1, ""},
	PersonRelation{"Napoleon V Victor", "Napoleon VI Louis", "child", 1, ""},
	PersonRelation{"Napoleon VI Louis", "Napoleon VII Charles", "child", 1, ""},
	PersonRelation{"Napoleon VII Charles", "Napoleon VIII Jean-Christophe", "child", 1, ""},
	PersonRelation{"Napoleon VII Charles", "Sophie Catherine Bonaparte", "child", 2, ""},
}

const (
	ancestryUrl = "tcp://scott@localhost:8222"
	//url1 = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
	//url2 = "tcp://scott@localhost:8222/{connectTimeout=30}";
	//url3 = "tcp://scott@localhost:8222/{dbName=mod;verifyDBName=true}";
	//url4 = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
	//url5 = "ssl://scott@localhost:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
	ancestryUser     = "napoleon"
	ancestryPassword = "bonaparte"
	//prefetchMetaData = false
	startYear = 1760
	endYear   = 1770
)

//Insert node data into database
func insertAncestryNodes(conn types.TGConnection, gof types.TGGraphObjectFactory) map[string]types.TGNode {
	fmt.Println(">>>>>>> Entering InsertAncestryNodes: Insert Few Family Nodes with individual properties <<<<<<<")

	var houseMemberTable = make(map[string]types.TGNode, 0)

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from InsertAncestryNodes - error during conn.GetGraphMetadata <<<<<<<")
		return nil
	}

	testNodeType, err := gmd.GetNodeType("houseMemberType")
	if err != nil {
		fmt.Println(">>>>>>> Returning from InsertAncestryNodes - error during conn.GetNodeType('houseMemberType') <<<<<<<")
		return nil
	}
	if testNodeType != nil {
		fmt.Printf(">>>>>>> 'houseMemberType' is found with %d attributes <<<<<<<\n", len(testNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'houseMemberType' is not found from meta data fetch <<<<<<<")
		return nil
	}

	for _, houseMember := range HouseMemberData {
		node1, err := gof.CreateNodeInGraph(testNodeType)
		if err != nil {
			fmt.Println(">>>>>>> Returning from InsertAncestryNodes - error during gof.CreateNode(node1) <<<<<<<")
			return nil
		}
		_ = node1.SetOrCreateAttribute("memberName", houseMember.MemberName)
		_ = node1.SetOrCreateAttribute("crownName", houseMember.CrownName)
		_ = node1.SetOrCreateAttribute("houseHead", houseMember.HouseHead)
		_ = node1.SetOrCreateAttribute("yearBorn", houseMember.YearBorn)
		_ = node1.SetOrCreateAttribute("yearDied", houseMember.YearDied)
		_ = node1.SetOrCreateAttribute("crownTitle", houseMember.CrownTitle)

		if houseMember.ReignStart != "" {
			reignStart, _ := time.Parse("02 Jan 2006", houseMember.ReignStart)
			_ = node1.SetOrCreateAttribute("reignStart", reignStart)
		} else {
			_ = node1.SetOrCreateAttribute("reignStart", nil)
		}

		if houseMember.ReignStart != "" {
			reignEnd, _ := time.Parse("02 Jan 2006", houseMember.ReignEnd)
			_ = node1.SetOrCreateAttribute("reignEnd", reignEnd)
		} else {
			_ = node1.SetOrCreateAttribute("reignEnd", nil)
		}

		err = conn.InsertEntity(node1)
		if err != nil {
			fmt.Println(">>>>>>> Returning from InsertAncestryNodes w/ error during conn.InsertEntity(node1) <<<<<<<")
			return nil
		}

		_, err = conn.Commit()
		if err != nil {
			fmt.Println(">>>>>>> Returning from InsertAncestryNodes w/ error during conn.Commit() <<<<<<<")
			return nil
		}
		fmt.Printf(">>>>>>> Inside InsertAncestryNodes: Successfully added node '%+v'<<<<<<<\n", houseMember.MemberName)
		houseMemberTable[houseMember.MemberName] = node1
	} // End of for loop
	fmt.Println(">>>>>>> Successfully added nodes w/ NO ERRORS !!! <<<<<<<")

	fmt.Println(">>>>>>> Returning from InsertAncestryNodes w/ NO ERRORS !!! <<<<<<<")
	return houseMemberTable
}

//Insert edge data into database
func insertRelationEdges(conn types.TGConnection, gof types.TGGraphObjectFactory, houseMemberTable map[string]types.TGNode) {
	fmt.Println(">>>>>>> Entering InsertRelationEdges: Insert Few Family Relations with individual properties <<<<<<<")

	// Insert edge data into database
	// Two edge types defined in ancestry-initdb.conf.
	// Added year of marriage and place of marriage edge attributes for spouseEdge desc
	// Added Birth order edge attribute for offspringEdge desc

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from InsertRelationEdges - error during conn.GetGraphMetadata <<<<<<<")
		return
	}

	spouseEdgeType, err := gmd.GetEdgeType("spouseEdge")
	if err != nil {
		fmt.Println(">>>>>>> Returning from InsertRelationEdges - error during conn.GetEdgeType('spouseEdge') <<<<<<<")
		return
	}
	if spouseEdgeType != nil {
		fmt.Printf(">>>>>>> 'spouseEdge' is found with %d attributes <<<<<<<\n", len(spouseEdgeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'spouseEdge' is not found from meta data fetch <<<<<<<")
		return
	}

	offspringEdgeType, err := gmd.GetEdgeType("offspringEdge")
	if err != nil {
		fmt.Println(">>>>>>> Returning from InsertRelationEdges - error during conn.GetEdgeType('offspringEdge') <<<<<<<")
		return
	}
	if offspringEdgeType != nil {
		fmt.Printf(">>>>>>> 'offspringEdgeType' is found with %d attributes <<<<<<<\n", len(offspringEdgeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'spouseEdge' is not found from meta data fetch <<<<<<<")
		return
	}

	for _, houseRelation := range HouseRelationData {
		houseMemberFrom := houseMemberTable[houseRelation.FromMemberName]
		houseMemberTo := houseMemberTable[houseRelation.ToMemberName]
		relationName := houseRelation.RelationDesc
		fmt.Printf(">>>>>>> Inside InsertRelationEdges: trying to create edge('%s'): From '%s' To '%s' <<<<<<<\n", relationName, houseRelation.FromMemberName, houseRelation.ToMemberName)
		//var relationDirection types.TGDirectionType
		if relationName == "spouse" {
			//relationDirection = types.DirectionTypeUnDirected
			spouseEdgeType.GetFromNodeType()
			edge1, err := gof.CreateEdgeWithEdgeType(houseMemberFrom, houseMemberTo, spouseEdgeType)
			if err != nil {
				fmt.Println(">>>>>>> Returning from InsertRelationEdges - error during gof.CreateEdgeWithEdgeType(spouseEdgeType) <<<<<<<")
				return
			}
			_ = edge1.SetOrCreateAttribute("yearMarried", houseRelation.Attribute1)
			_ = edge1.SetOrCreateAttribute("placeMarried", houseRelation.Attribute2)
			_ = edge1.SetOrCreateAttribute("relType", relationName)
			err = conn.InsertEntity(edge1)
			if err != nil {
				fmt.Println(">>>>>>> Returning from InsertRelationEdges - error during conn.InsertEntity(edge1) <<<<<<<")
				return
			}

			_, err = conn.Commit()
			if err != nil {
				fmt.Println(">>>>>>> Returning from InsertRelationEdges w/ error during conn.Commit() <<<<<<<")
				return
			}
			fmt.Printf(">>>>>>> Inside InsertRelationEdges: Successfully added edge(spouse): From '%+v' To '%+v' <<<<<<<\n", houseRelation.FromMemberName, houseRelation.ToMemberName)
		} else {
			//relationDirection = types.DirectionTypeDirected
			edge1, err := gof.CreateEdgeWithEdgeType(houseMemberFrom, houseMemberTo, offspringEdgeType)
			if err != nil {
				fmt.Println(">>>>>>> Returning from InsertRelationEdges - error during gof.CreateEdgeWithEdgeType(offspringEdgeType) <<<<<<<")
				return
			}
			_ = edge1.SetOrCreateAttribute("birthOrder", houseRelation.Attribute1)
			_ = edge1.SetOrCreateAttribute("relType", relationName)
			err = conn.InsertEntity(edge1)
			if err != nil {
				fmt.Println(">>>>>>> Returning from InsertRelationEdges - error during conn.InsertEntity(edge1) <<<<<<<")
				return
			}

			_, err = conn.Commit()
			if err != nil {
				fmt.Println(">>>>>>> Returning from InsertRelationEdges w/ error during conn.Commit() <<<<<<<")
				return
			}
			fmt.Printf(">>>>>>> Inside InsertRelationEdges: Successfully added edge: From '%+v' To '%+v' <<<<<<<\n", houseRelation.FromMemberName, houseRelation.ToMemberName)
		}
	} // End of for loop

	fmt.Println(">>>>>>> Successfully added edges w/ NO ERRORS !!! <<<<<<<")
	fmt.Println(">>>>>>> Returning from InsertRelationEdges w/ NO ERRORS !!! <<<<<<<")
}

/**
 * Modified from example BuildGraph to demonstrate traversal filter
 * Uses ancestry-initdb.conf and ancestry-tgdb.conf
 *
 * Build the House of Bonaparte graph.
 * The House of Bonaparte is an imperial and royal European dynasty founded in 1804 by Napoleon I.
 *
 * Each member has the following characteristics :
 * -memberName : Name of member (primary key)
 * -crownName : Name while reigning
 * -crownTitle : Title while reigning
 * -houseHead : There is always a head of the house at a given time.
 * -yearBorn : Year of birth
 * -yearDied : Year of death
 * -reignStart : Date reign started
 * -reignEnd : Date reign ended
 *
 */
//Build Ancestry Graph / Tree based on the available data
func BuildAncestryGraph() {
	fmt.Println(">>>>>>> Entering BuildAncestryGraph <<<<<<<")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(ancestryUrl, ancestryUser, ancestryPassword, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from BuildAncestryGraph - error during CreateConnection <<<<<<<")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from BuildAncestryGraph - error during conn.Connect <<<<<<<")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from BuildAncestryGraph - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from BuildAncestryGraph - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside BuildAncestryGraph: About to InsertAncestryNodes <<<<<<<")
	houseMemberTable := insertAncestryNodes(conn, gof)
	fmt.Println(">>>>>>> Inside BuildAncestryGraph: About to InsertRelationEdges <<<<<<<")
	insertRelationEdges(conn, gof, houseMemberTable)
	fmt.Println(">>>>>>> Inside BuildAncestryGraph: Napoleon Bonaparte Ancestry Graph Created Successfully <<<<<<<")

	err = conn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from BuildAncestryGraph - error during conn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from BuildAncestryGraph - successfully disconnected. <<<<<<<")
}

func executeQuery1(conn types.TGConnection, startYear, endYear int) {
	fmt.Println(">>>>>>> Entering executeQuery1: Query Family Relations based on the conditions <<<<<<<")

	//Simple query
	//dumpDepth := 5
	//currDepth := 0
	//dumpBreadth := false
	//showAllPath := true
	//option := query.NewQueryOption()
	//dumpBreadth = true
	queryString1 := fmt.Sprintf("@nodetype = 'houseMemberType' and yearBorn > %d and yearBorn < %d;", startYear, endYear)
	fmt.Printf(">>>>>>> Inside executeQuery1: Executing %s <<<<<<<\n", queryString1)

	resultSet, err := conn.ExecuteQuery(queryString1, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from executeQuery1 - error during conn.ExecuteQuery('queryString1') <<<<<<<")
		return
	}
	if resultSet != nil {
		for {
			if !resultSet.HasNext() {
				break
			}
			houseMember := resultSet.Next()
			//if (dumpBreadth) {
			//	EntityUtils.printEntitiesBreadth((TGNode) houseMember, dumpDepth);
			//} else {
			//	EntityUtils.printEntities(houseMember, dumpDepth, currDepth, "", showAllPath, new HashMap<Integer, TGEntity>());
			//}
			fmt.Printf(">>>>>>> Query returned house member '%+v'<<<<<<<\n", houseMember)
		}
	} else {
		fmt.Printf(">>>>>>> Query '%s' DID NOT Return any results <<<<<<<\n", queryString1)
		return
	}

	fmt.Println(">>>>>>> Returning from executeQuery1 w/ NO ERRORS !!! <<<<<<<")
}

func executeQuery2(conn types.TGConnection) {
	fmt.Println(">>>>>>> Entering executeQuery2: Query Family Relations based on the conditions <<<<<<<")

	//Identify a single path from Napoleon Bonaparte  to Francois Bonaparte
	//dumpDepth := 5
	//currDepth := 0
	//dumpBreadth := false
	//showAllPath := true
	//option := query.NewQueryOption()
	//dumpBreadth = true
	queryString2 := fmt.Sprint("@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';")
	traverseString := fmt.Sprint("@edgetype = 'offspringEdge' and @isfromedge = 1 and @edge.birthOrder = 1 and @degree < 3;")
	endString := fmt.Sprint("@tonodetype = 'houseMemberType' and @tonode.memberName = 'Francois Bonaparte';")
	fmt.Printf(">>>>>>> Inside executeQuery2: Executing %s <<<<<<<\n", queryString2)

	resultSet, err := conn.ExecuteQueryWithFilter(queryString2, "", traverseString, endString, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from executeQuery2 - error during conn.GetEdgeType('queryString2') <<<<<<<")
		return
	}
	if resultSet != nil {
		for {
			if !resultSet.HasNext() {
				break
			}
			houseMember := resultSet.Next()
			//if (dumpBreadth) {
			//	EntityUtils.printEntitiesBreadth((TGNode) houseMember, dumpDepth);
			//} else {
			//	EntityUtils.printEntities(houseMember, dumpDepth, currDepth, "", showAllPath, new HashMap<Integer, TGEntity>());
			//}
			fmt.Printf(">>>>>>> Query returned house member '%+v'<<<<<<<\n", houseMember)
		}
	} else {
		fmt.Printf(">>>>>>> Query '%s' DID NOT Return any results <<<<<<<\n", queryString2)
		return
	}

	fmt.Println(">>>>>>> Returning from executeQuery2 w/ NO ERRORS !!! <<<<<<<")
}

func executeQuery3(conn types.TGConnection) {
	fmt.Println(">>>>>>> Entering executeQuery3: Query Family Relations based on the conditions <<<<<<<")

	//Identify all paths from Napoleon Bonaparte  to Napoleon IV Eugene with no traversal filter except 10 level deep restriction
	//dumpDepth := 10
	//currDepth := 0
	//dumpBreadth := false
	//showAllPath := true
	//option := query.NewQueryOption()
	//dumpBreadth = true
	queryString3 := fmt.Sprint("@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';")
	endString := fmt.Sprint("@tonodetype = 'houseMemberType' and @tonode.memberName = 'Napoleon IV Eugene';")
	fmt.Printf(">>>>>>> Inside executeQuery3: Executing %s <<<<<<<\n", queryString3)

	resultSet, err := conn.ExecuteQueryWithFilter(queryString3, "", "", endString, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from executeQuery3 - error during conn.GetEdgeType('spouseEdge') <<<<<<<")
		return
	}
	if resultSet != nil {
		for {
			if !resultSet.HasNext() {
				break
			}
			houseMember := resultSet.Next()
			//if (dumpBreadth) {
			//	EntityUtils.printEntitiesBreadth((TGNode) houseMember, dumpDepth);
			//} else {
			//	EntityUtils.printEntities(houseMember, dumpDepth, currDepth, "", showAllPath, new HashMap<Integer, TGEntity>());
			//}
			fmt.Printf(">>>>>>> Query returned house member '%+v'<<<<<<<\n", houseMember)
		}
	} else {
		fmt.Printf(">>>>>>> Query '%s' DID NOT Return any results <<<<<<<\n", queryString3)
		return
	}

	fmt.Println(">>>>>>> Returning from executeQuery3 w/ NO ERRORS !!! <<<<<<<")
}

func executeQuery4(conn types.TGConnection) {
	fmt.Println(">>>>>>> Entering executeQuery4: Query Family Relations based on the conditions <<<<<<<")

	//Identify all paths from Napoleon Bonaparte  to Napoleon IV Eugene using only offspringEdge desc and within 10 level deep
	//dumpDepth := 10
	//currDepth := 0
	//dumpBreadth := false
	//showAllPath := true
	option := query.NewQueryOption()
	_ = option.SetTraversalDepth(10)
	//dumpBreadth = true
	queryString4 := fmt.Sprint("@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';")
	traverseString := fmt.Sprint("@edgetype = 'offspringEdge' and @degree <= 10;")
	endString := fmt.Sprint("@tonodetype = 'houseMemberType' and @tonode.memberName = 'Napoleon IV Eugene';")
	fmt.Printf(">>>>>>> Inside executeQuery4: Executing %s <<<<<<<\n", queryString4)

	resultSet, err := conn.ExecuteQueryWithFilter(queryString4, "", traverseString, endString, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from executeQuery4 - error during conn.GetEdgeType('spouseEdge') <<<<<<<")
		return
	}
	if resultSet != nil {
		for {
			if !resultSet.HasNext() {
				break
			}
			houseMember := resultSet.Next()
			//if (dumpBreadth) {
			//	EntityUtils.printEntitiesBreadth((TGNode) houseMember, dumpDepth);
			//} else {
			//	EntityUtils.printEntities(houseMember, dumpDepth, currDepth, "", showAllPath, new HashMap<Integer, TGEntity>());
			//}
			fmt.Printf(">>>>>>> Query returned house member '%+v'<<<<<<<\n", houseMember)
		}
	} else {
		fmt.Printf(">>>>>>> Query '%s' DID NOT Return any results <<<<<<<\n", queryString4)
		return
	}

	fmt.Println(">>>>>>> Returning from executeQuery4 w/ NO ERRORS !!! <<<<<<<")
}

func executeQuery5(conn types.TGConnection) {
	fmt.Println(">>>>>>> Entering executeQuery5: Query Family Relations based on the conditions <<<<<<<")

	//Identify specific path from Napoleon Bonaparte -> his parents -> Louis Bonaparte -> Louis Napoleon -> Napoleon IV Eugene
	//dumpDepth := 10
	//currDepth := 0
	//dumpBreadth := false
	//showAllPath := true
	option := query.NewQueryOption()
	_ = option.SetTraversalDepth(10)
	//dumpBreadth = true
	queryString5 := fmt.Sprint("@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';")
	traverseString := fmt.Sprint("(@edgetype = 'offspringEdge' and @isfromedge = 0 and @degree = 1) or (@edgetype = 'offspringEdge' and @isfromedge = 1 and @degree = 2) or (@edgetype = 'offspringEdge' and @isfromedge = 1 and @degree = 3) or (@edgetype = 'offspringEdge' and @isfromedge = 1 and @degree = 4);")
	endString := fmt.Sprint("@tonodetype = 'houseMemberType' and @tonode.memberName = 'Napoleon IV Eugene';")
	fmt.Printf(">>>>>>> Inside executeQuery5: Executing %s <<<<<<<\n", queryString5)

	resultSet, err := conn.ExecuteQueryWithFilter(queryString5, "", traverseString, endString, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from executeQuery5 - error during conn.GetEdgeType('spouseEdge') <<<<<<<")
		return
	}
	if resultSet != nil {
		for {
			if !resultSet.HasNext() {
				break
			}
			houseMember := resultSet.Next()
			//if (dumpBreadth) {
			//	EntityUtils.printEntitiesBreadth((TGNode) houseMember, dumpDepth);
			//} else {
			//	EntityUtils.printEntities(houseMember, dumpDepth, currDepth, "", showAllPath, new HashMap<Integer, TGEntity>());
			//}
			fmt.Printf(">>>>>>> Query returned house member '%+v'<<<<<<<\n", houseMember)
		}
	} else {
		fmt.Printf(">>>>>>> Query '%s' DID NOT Return any results <<<<<<<\n", queryString5)
		return
	}

	fmt.Println(">>>>>>> Returning from executeQuery5 w/ NO ERRORS !!! <<<<<<<")
}

/**
 *
 * Modified from example QueryGraph to demonstrate traversal filter
 * Uses ancestry-initdb.conf and ancestry-tgdb.conf
 * Needs to execute BuildAncestryGraph() first
 *
 * Traversal starts from the nodes in the query results set.
 * Starting from each edge of those nodes the traversal filter is applied.
 * If evaluation is true, it will apply the end filter to see the ending
 * condition is fulfilled.  If not, it will continue the traversal to the
 * next level. The final result set only contains the nodes and edges
 * fulfilled all the conditions.
 *
 * The following reserved keyword are introduced for traversal filtering
 * @fromnodetype - string - the node desc name of the node we are getting the edge from
 * @tonodetype - string - the node desc name of the other end of the edge
 * @isfromedge - is the node where the edge is retrieve from on the from side of the edge. 1 - true, 0 - false
 * @fromnode.<attr name> - retrieve the from node attribute value
 * @tonode.<attr name> - retrieve the to node attribute value
 * @edge.<attr name> - retrieve the edge attribute value
 * @edgetype - string - edge desc name
 * @degree/depth - int - degree of separation or what we call the depth
 *                       both degree and depth are valid keywords
 *
 * e.g.  If we starts from Napoleon Bonaparte and get the offspring edge to Carlo Bonaparte,
 *       the isfromedge will be 0 because the edge is created from Carlo to Napoleon.
 *       But the fromnode is Napoleon because we are traversing from Napoleon to Carlo
 *       and therefore tonode is Carlo.  It may sound confusing.  May be a different
 *       naming for fromnode and tonode can help.
 *
 * The isfromedge is used to control which direction of the edge you want to traverse.
 * Using offspring edge as an example, if you only interested in traversing from
 * parent to child, you should specify isfromedge to 1. If you start from
 * Napoleon Bonaparte and traverse offspring edge with isfromedge = 1, you will
 * only get to Francois Bonaparte but not to Carlo or Letitia.
 *
 * The second argument of the new executeQuery method is not used right now.  The end
 * condition is required but the traversal condition is optional.
 *
 * Query for members in the House of Bonaparte graph
 * born between the start and end years
 * and display the member attributes.
 *
 */
//Query Ancestry Graph / Tree based on the available inputs
func QueryAncestryGraph() {
	fmt.Println(">>>>>>> Entering QueryAncestryGraph <<<<<<<")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(ancestryUrl, ancestryUser, ancestryPassword, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from QueryAncestryGraph - error during CreateConnection <<<<<<<")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from QueryAncestryGraph - error during conn.Connect <<<<<<<")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from QueryAncestryGraph - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from QueryAncestryGraph - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside QueryAncestryGraph: About to executeQuery1 <<<<<<<")
	executeQuery1(conn, startYear, endYear)
	fmt.Println(">>>>>>> Inside QueryAncestryGraph: About to executeQuery2 <<<<<<<")
	executeQuery2(conn)
	fmt.Println(">>>>>>> Inside QueryAncestryGraph: About to executeQuery3 <<<<<<<")
	executeQuery3(conn)
	fmt.Println(">>>>>>> Inside QueryAncestryGraph: About to executeQuery4 <<<<<<<")
	executeQuery4(conn)
	fmt.Println(">>>>>>> Inside QueryAncestryGraph: About to executeQuery5 <<<<<<<")
	executeQuery5(conn)
	fmt.Println(">>>>>>> Inside QueryAncestryGraph: Napoleon Bonaparte Ancestry Graph Queries Executed Successfully <<<<<<<")

	err = conn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from QueryAncestryGraph - error during conn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from QueryAncestryGraph - successfully disconnected. <<<<<<<")
}

func AncestryTest() {
	fmt.Println(">>>>>>> Entering BuildAncestryGraph <<<<<<<")

	fmt.Println(">>>>>>> Inside BuildAncestryGraph: About to BuildAncestryGraph <<<<<<<")
	BuildAncestryGraph()
	fmt.Println(">>>>>>> Inside BuildAncestryGraph: About to QueryAncestryGraph <<<<<<<")
	QueryAncestryGraph()

	fmt.Println(">>>>>>> Returning from BuildAncestryGraph - successfully disconnected. <<<<<<<")
}
