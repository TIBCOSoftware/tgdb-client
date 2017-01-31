/**
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
 * Usage : node BuildGraph.js
 * 
 */	
var TGException  = require('../lib/exception/TGException').TGException,
    conFactory   = require('../lib/connection/TGConnectionFactory'),
	PrintUtility = require('../lib/utils/PrintUtility'),
	TGLogManager = require('../lib/log/TGLogManager'),
	TGLogLevel   = require('../lib/log/TGLogger').TGLogLevel;
	
var logger = TGLogManager.getLogger();

// Members of the House to be inserted in the database
var houseMemberData = [
	// memberName, crownName, houseHead, yearBorn, yearDied, reignStart, reignEnd, crownTitle
	[ "Carlo Bonaparte", null, false, 1746, 1785, null, null, null ],
	[ "Letizia Ramolino", null, false, 1750, 1836, null, null, null ],
	[ "Joseph Bonaparte", "Joseph I", false, 1768, 1844, "6 Jun 1808", "11 Dec 1813", "King of Spain" ],
	[ "Napoleon Bonaparte", "Napoleon I", false, 1769, 1821, "18 May 1804", "22 Jun 1815", "Emperor of the French" ],
	[ "Lucien Bonaparte", null, false, 1775, 1840, null, null, null ],
	[ "Elisa Bonaparte", "Elisa Bonaparte", false, 1777, 1820, "3 Mar 1809", "1 Feb 1814", "Grand Duchess of Tuscany" ],
	[ "Louis Bonaparte", "Louis I", false, 1778, 1846, "5 Jun 1806", "1 Jul 1810", "King of Holland" ],
	[ "Pauline Bonaparte", null, false, 1780, 1825, null, null, null ],
	[ "Caroline Bonaparte", null, false, 1782, 1839, null, null, null ],
	[ "Jerome Bonaparte", "Jerome I", false, 1784, 1860, "8 Jul 1807", "26 Oct 1813", "King of Westphalia" ],
	[ "Marie Louise of Austria", null, false, 1791, 1847, null, null, "Empress Consort of the French" ],
	[ "Josephine of Beauharnais", null, false, 1763, 1814, null, null, "Empress Consort of the French" ],
	[ "Alexandre of Beauharnais", null, false, 1760, 1794, null, null, null ],
	[ "Betsy Patterson", null, false, 1785, 1879, null, null, null ],
	[ "Catharina of Wurttemberg", null, false, 1783, 1835, null, null, "Queen Consort of Westphalia" ],
	[ "Francois Bonaparte", "Napoleon II", false, 1811, 1832, "22 Jun 1815", "7 Jul 1815", "Emperor of the French" ],
	[ "Hortense of Beauharnais", null, false, 1783, 1837, null, null, "Queen Consort of Holland" ],
	[ "Jerome Napoleon", null, false, 1805, 1870, null, null, null ],
	[ "Prince Napoleon", null, false, 1822, 1891, null, null, null ],
	[ "Louis Napoleon", "Napoleon III", true, 1808, 1873, "2 Dec 1852", "4 Sep 1870", "Emperors of the French" ],
	[ "Napoleon-Louis Bonaparte", "Louis II", false, 1804, 1831, "1 Jul 1810", "13 Jul 1810", "King of Holland" ],
	[ "Napoleon IV Eugene", null, true, 1856, 1879, null, null, null ],
	[ "Napoleon V Victor", null, true, 1862, 1926, null, null, null ],
	[ "Marie Clotilde Bonaparte", null, false, 1912, 1996, null, null, null ],
	[ "Napoleon VI Louis", null, true, 1914, 1997, null, null, null ],
	[ "Napoleon VII Charles", null, true, 1950, null, null, null, null ],
	[ "Napoleon VIII Jean-Christophe", null, true, 1986, null, null, null, null ],
	[ "Sophie Catherine Bonaparte", null, false, 1992, null, null, null, null ] 
];

// Relation among the members of the House to be inserted in the database
var houseRelationData2 = [
// From memberName, To memberName, relation type
	[ "Carlo Bonaparte", "Letizia Ramolino", "spouse" ], 
	[ "Joseph Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Joseph Bonaparte", "Letizia Ramolino", "parent" ], 
	[ "Napoleon Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Napoleon Bonaparte", "Letizia Ramolino", "parent" ], 
	[ "Lucien Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Lucien Bonaparte", "Letizia Ramolino", "parent" ], 
	[ "Elisa Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Elisa Bonaparte", "Letizia Ramolino", "parent" ], 
	[ "Louis Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Louis Bonaparte", "Letizia Ramolino", "parent" ], 
	[ "Pauline Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Pauline Bonaparte", "Letizia Ramolino", "parent" ], 
	[ "Caroline Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Caroline Bonaparte", "Letizia Ramolino", "parent" ], 
	[ "Jerome Bonaparte", "Carlo Bonaparte", "parent" ],
	[ "Jerome Bonaparte", "Letizia Ramolino", "parent" ],

	[ "Napoleon Bonaparte", "Marie Louise of Austria", "spouse" ],
	[ "Francois Bonaparte", "Napoleon Bonaparte", "parent" ],
	[ "Francois Bonaparte", "Marie Louise of Austria", "parent" ],

	[ "Napoleon Bonaparte", "Josephine of Beauharnais", "spouse" ],

	[ "Alexandre of Beauharnais", "Josephine of Beauharnais", "spouse" ],
	[ "Hortense of Beauharnais", "Alexandre of Beauharnais", "parent" ],
	[ "Hortense of Beauharnais", "Josephine of Beauharnais", "parent" ],

	[ "Louis Bonaparte", "Hortense of Beauharnais", "spouse" ],
	[ "Louis Napoleon", "Louis Bonaparte", "parent" ],
	[ "Louis Napoleon", "Hortense of Beauharnais", "parent" ],
	[ "Napoleon-Louis Bonaparte", "Louis Bonaparte", "parent" ],
	[ "Napoleon-Louis Bonaparte", "Hortense of Beauharnais", "parent" ],

	[ "Jerome Bonaparte", "Betsy Patterson", "spouse" ],
	[ "Jerome Napoleon", "Jerome Bonaparte", "parent" ],
	[ "Jerome Napoleon", "Betsy Patterson", "parent" ],

	[ "Jerome Bonaparte", "Catharina of Wurttemberg", "spouse" ],
	[ "Prince Napoleon", "Jerome Bonaparte", "parent" ],
	[ "Prince Napoleon", "Catharina of Wurttemberg", "parent" ],

	[ "Napoleon IV Eugene", "Louis Napoleon", "parent" ],

	[ "Napoleon V Victor", "Prince Napoleon", "parent" ],

	[ "Napoleon VI Louis", "Napoleon V Victor", "parent" ],

	[ "Napoleon VII Charles", "Napoleon VI Louis", "parent" ],

	[ "Napoleon VIII Jean-Christophe", "Napoleon VII Charles", "parent" ],
	[ "Sophie Catherine Bonaparte", "Napoleon VII Charles", "parent" ] 
];

//Relation among the members of the House to be inserted in the database
var houseRelationData = [
// From memberName, To memberName, relation type
	[ "Carlo Bonaparte", "Letizia Ramolino", "spouse" ], 
	[ "Carlo Bonaparte", "Joseph Bonaparte", "child" ],
	[ "Letizia Ramolino", "Joseph Bonaparte", "child" ], 
	[ "Carlo Bonaparte", "Napoleon Bonaparte", "child" ],
	[ "Letizia Ramolino", "Napoleon Bonaparte", "child" ], 
	[ "Carlo Bonaparte", "Lucien Bonaparte", "child" ],
	[ "Letizia Ramolino", "Lucien Bonaparte", "child" ], 
	[ "Carlo Bonaparte", "Elisa Bonaparte", "child" ],
	[ "Letizia Ramolino", "Elisa Bonaparte", "child" ], 
	[ "Carlo Bonaparte", "Louis Bonaparte", "child" ],
	[ "Letizia Ramolino", "Louis Bonaparte", "child" ], 
	[ "Carlo Bonaparte", "Pauline Bonaparte", "child" ],
	[ "Letizia Ramolino", "Pauline Bonaparte", "child" ], 
	[ "Carlo Bonaparte", "Caroline Bonaparte", "child" ],
	[ "Letizia Ramolino", "Caroline Bonaparte", "child" ], 
	[ "Carlo Bonaparte", "Jerome Bonaparte", "child" ],
	[ "Letizia Ramolino", "Jerome Bonaparte", "child" ],

	[ "Napoleon Bonaparte", "Marie Louise of Austria", "spouse" ],
	[ "Napoleon Bonaparte", "Francois Bonaparte", "child" ],
	[ "Marie Louise of Austria", "Francois Bonaparte", "child" ],

	[ "Napoleon Bonaparte", "Josephine of Beauharnais", "spouse" ],

	[ "Alexandre of Beauharnais", "Josephine of Beauharnais", "spouse" ],
	[ "Alexandre of Beauharnais", "Hortense of Beauharnais", "child" ],
	[ "Josephine of Beauharnais", "Hortense of Beauharnais", "child" ],

	[ "Louis Bonaparte", "Hortense of Beauharnais", "spouse" ],
	[ "Louis Bonaparte", "Louis Napoleon", "child" ],
	[ "Hortense of Beauharnais", "Louis Napoleon", "child" ],
	[ "Louis Bonaparte", "Napoleon-Louis Bonaparte", "child" ],
	[ "Hortense of Beauharnais", "Napoleon-Louis Bonaparte", "child" ],

	[ "Jerome Bonaparte", "Betsy Patterson", "spouse" ],
	[ "Jerome Bonaparte", "Jerome Napoleon", "child" ],
	[ "Betsy Patterson", "Jerome Napoleon", "child" ],

	[ "Jerome Bonaparte", "Catharina of Wurttemberg", "spouse" ],
	[ "Jerome Bonaparte", "Prince Napoleon", "child" ],
	[ "Catharina of Wurttemberg", "Prince Napoleon", "child" ],

	[ "Louis Napoleon", "Napoleon IV Eugene", "child" ],

	[ "Prince Napoleon", "Napoleon V Victor", "child" ],

	[ "Napoleon V Victor", "Napoleon VI Louis", "child" ],

	[ "Napoleon VI Louis", "Napoleon VII Charles", "child" ],

	[ "Napoleon VII Charles", "Napoleon VIII Jean-Christophe", "child" ],
	[ "Napoleon VII Charles", "Sophie Catherine Bonaparte", "child" ] 
];


var houseMemberTable = {};
function insertHouseMemberData(conn, gof, houseMemberType, houseMemberData, callback) {
	if(houseMemberData.length>0) {
		var dataRow = houseMemberData.shift();
		var houseMember = gof.createNode(houseMemberType);
		houseMember.setAttribute("memberName", dataRow[0]);
		houseMember.setAttribute("crownName", dataRow[1]);
		houseMember.setAttribute("houseHead", dataRow[2]);
		houseMember.setAttribute("yearBorn", dataRow[3]);
		houseMember.setAttribute("yearDied", dataRow[4]);
		houseMember.setAttribute("crownTitle", dataRow[7]);
		
		if (dataRow[5] !== null) {
			houseMember.setAttribute("reignStart", new Date(dataRow[5]));
		}
		else {
			houseMember.setAttribute("reignStart", null);
		}
		
		if (dataRow[6] !== null) {
			houseMember.setAttribute("reignEnd", new Date(dataRow[6]));
		}
		else {
			houseMember.setAttribute("reignEnd", null);
		}
		
		logger.logInfo(
				"Transaction started for Node : " + 
				houseMember.getAttribute("memberName").getValue());
		console.log('Process house member data row : %s %s %s %s %s %s', 
				dataRow[0], dataRow[1], dataRow[2], 
				dataRow[3], dataRow[4], dataRow[7]);
		conn.insertEntity(houseMember);
		conn.commit(function(){
			logger.logInfo(
					"Transaction completed for Node : " + 
					houseMember.getAttribute("memberName").getValue());
			houseMemberTable[houseMember.getAttribute("memberName").getValue()] = houseMember;
			logger.logDebug(
					"Store node for key = %s, val = %s", 
					houseMember.getAttribute("memberName").getValue(), houseMember);
			insertHouseMemberData(conn, gof, houseMemberType, houseMemberData, callback);
		}); // Write data to database
	}
	else {
		callback(); // return;
	}
}

function insertHouseRelationData(conn, gof, houseRelationData, callback) {	
	if(houseRelationData.length>0) {
		var dataRow = houseRelationData.shift();
		var houseMemberFrom = houseMemberTable[dataRow[0]];
		logger.logDebug("Lookup from node for key = %s, val = %s", dataRow[0], houseMemberTable[dataRow[0]]);

		var houseMemberTo = houseMemberTable[dataRow[1]];
		logger.logDebug("Lookup to node for key = %s, val = %s", dataRow[1], houseMemberTable[dataRow[1]]);
		
		var houseRelation = null;
		if(dataRow[2] === "spouse") {
			houseRelation = gof.createUndirectedEdge(houseMemberFrom, houseMemberTo);
		} else {
			houseRelation = gof.createDirectedEdge(houseMemberFrom, houseMemberTo);
		}
		houseRelation.setAttribute("relType", dataRow[2]);
		
		logger.logInfo(
				"Transaction Started for Edge : " + 
				houseMemberFrom.getAttribute("memberName").getValue() + 
				" to " + houseMemberTo.getAttribute("memberName").getValue());
		console.log('Process house relation data row : %s %s %s', 
				dataRow[0], dataRow[1], dataRow[2]);
		conn.insertEntity(houseRelation);
		conn.commit(function(){
			logger.logInfo(
					"Transaction completed for Edge : " + 
					houseMemberFrom.getAttribute("memberName").getValue() + 
					" to " + houseMemberTo.getAttribute("memberName").getValue());
			insertHouseRelationData(conn, gof, houseRelationData, callback);
		});
	} else {
		callback(); // return
	}
}

function main() {
	var logger = TGLogManager.getLogger();
	logger.setLevel(TGLogLevel.Info);	
	
	var url = "tcp://192.168.1.6:8222";
	var user = "napoleon";
	var pwd = "bonaparte";
	var conn = conFactory.getFactory().createConnection(url, user, pwd, null);
	conn.connect(function(){
		var gof = conn.getGraphObjectFactory();
		if (!gof) {
			throw new TGException("Graph object not found");
		}

		conn.getGraphMetadata(true, function(gmd){
			var houseMemberType = gmd.getNodeType("houseMemberType");
			if (!houseMemberType) {
				throw new TGException("Node type not found");
			}
			
			insertHouseMemberData(conn, gof, houseMemberType, houseMemberData, function(){
				logger.logInfo("-------------------------------------------------");
				insertHouseRelationData(conn, gof, houseRelationData, function(){
					logger.logInfo("\nHouse of Bonaparte graph completed successfully");
					if (conn !== null) {
						conn.disconnect();
					}
				});
			});
		});
	});
}

main();