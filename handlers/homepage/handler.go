package homepage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/onsdigital/dp-frontend-router/lang"
	"github.com/onsdigital/go-ns/log"
)

var stubbedData = `{
  "taxonomy": [
    {
      "uri": "/businessindustryandtrade",
      "children": [
        {
          "uri": "/businessindustryandtrade/business",
          "description": {},
          "children": [
            {
              "uri": "/businessindustryandtrade/business/activitysizeandlocation",
              "description": {
                "title": "Activity, size and location"
              },
              "type": "product_page"
            },
            {
              "uri": "/businessindustryandtrade/business/businessinnovation",
              "description": {
                "title": "Business innovation"
              },
              "type": "product_page"
            },
            {
              "uri": "/businessindustryandtrade/business/businessservices",
              "description": {
                "title": "Business services"
              },
              "type": "product_page"
            }
          ],
          "title": "Business"
        },
        {
          "uri": "/businessindustryandtrade/changestobusiness",
          "description": {},
          "children": [
            {
              "uri": "/businessindustryandtrade/changestobusiness/bankruptcyinsolvency",
              "description": {
                "title": "Bankruptcy/insolvency"
              },
              "type": "product_page"
            },
            {
              "uri": "/businessindustryandtrade/changestobusiness/businessbirthsdeathsandsurvivalrates",
              "description": {
                "title": "Business births, deaths and survival rates"
              },
              "type": "product_page"
            },
            {
              "uri": "/businessindustryandtrade/changestobusiness/mergersandacquisitions",
              "description": {
                "title": "Mergers and acquisitions"
              },
              "type": "product_page"
            }
          ],
          "title": "Changes to business"
        },
        {
          "uri": "/businessindustryandtrade/constructionindustry",
          "description": {},
          "title": "Construction industry"
        },
        {
          "uri": "/businessindustryandtrade/itandinternetindustry",
          "description": {},
          "title": "IT and internet industry"
        },
        {
          "uri": "/businessindustryandtrade/internationaltrade",
          "description": {},
          "title": "International trade"
        },
        {
          "uri": "/businessindustryandtrade/manufacturingandproductionindustry",
          "description": {},
          "title": "Manufacturing and production industry"
        },
        {
          "uri": "/businessindustryandtrade/retailindustry",
          "description": {},
          "title": "Retail industry"
        },
        {
          "uri": "/businessindustryandtrade/tourismindustry",
          "description": {},
          "title": "Tourism industry"
        }
      ],
      "title": "Business, industry and trade"
    },
    {
      "uri": "/economy",
      "children": [
        {
          "uri": "/economy/economicoutputandproductivity",
          "description": {},
          "children": [
            {
              "uri": "/economy/economicoutputandproductivity/output",
              "description": {
                "title": "Output"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/economicoutputandproductivity/productivitymeasures",
              "description": {
                "title": "Productivity measures"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/economicoutputandproductivity/publicservicesproductivity",
              "description": {
                "title": "Public services productivity"
              },
              "type": "product_page"
            }
          ],
          "title": "Economic output and productivity"
        },
        {
          "uri": "/economy/environmentalaccounts",
          "description": {},
          "title": "Environmental accounts"
        },
        {
          "uri": "/economy/governmentpublicsectorandtaxes",
          "description": {},
          "children": [
            {
              "uri": "/economy/governmentpublicsectorandtaxes/localgovernmentfinance",
              "description": {
                "title": "Local government finance"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/governmentpublicsectorandtaxes/publicsectorfinance",
              "description": {
                "title": "Public sector finance"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/governmentpublicsectorandtaxes/publicspending",
              "description": {
                "title": "Public spending"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/governmentpublicsectorandtaxes/researchanddevelopmentexpenditure",
              "description": {
                "title": "Research and development expenditure"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/governmentpublicsectorandtaxes/taxesandrevenue",
              "description": {
                "title": "Taxes and revenue"
              },
              "type": "product_page"
            }
          ],
          "title": "Government, public sector and taxes"
        },
        {
          "uri": "/economy/grossdomesticproductgdp",
          "description": {},
          "title": "Gross Domestic Product (GDP)"
        },
        {
          "uri": "/economy/grossvalueaddedgva",
          "description": {},
          "title": "Gross Value Added (GVA)"
        },
        {
          "uri": "/economy/inflationandpriceindices",
          "description": {},
          "title": "Inflation and price indices"
        },
        {
          "uri": "/economy/investmentspensionsandtrusts",
          "description": {},
          "title": "Investments, pensions and trusts"
        },
        {
          "uri": "/economy/nationalaccounts",
          "description": {},
          "children": [
            {
              "uri": "/economy/nationalaccounts/balanceofpayments",
              "description": {
                "title": "Balance of payments"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/nationalaccounts/satelliteaccounts",
              "description": {
                "title": "Satellite accounts"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/nationalaccounts/supplyandusetables",
              "description": {
                "title": "Supply and use tables"
              },
              "type": "product_page"
            },
            {
              "uri": "/economy/nationalaccounts/uksectoraccounts",
              "description": {
                "title": "UK sector accounts"
              },
              "type": "product_page"
            }
          ],
          "title": "National accounts"
        },
        {
          "uri": "/economy/regionalaccounts",
          "description": {},
          "children": [
            {
              "uri": "/economy/regionalaccounts/grossdisposablehouseholdincome",
              "description": {
                "title": "Gross disposable household income"
              },
              "type": "product_page"
            }
          ],
          "title": "Regional accounts"
        }
      ],
      "title": "Economy"
    },
    {
      "uri": "/employmentandlabourmarket",
      "children": [
        {
          "uri": "/employmentandlabourmarket/peopleinwork",
          "description": {},
          "children": [
            {
              "uri": "/employmentandlabourmarket/peopleinwork/earningsandworkinghours",
              "description": {
                "title": "Earnings and working hours"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peopleinwork/employmentandemployeetypes",
              "description": {
                "title": "Employment and employee types"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peopleinwork/labourproductivity",
              "description": {
                "title": "Labour productivity"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peopleinwork/publicsectorpersonnel",
              "description": {
                "title": "Public sector personnel"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peopleinwork/workplacedisputesandworkingconditions",
              "description": {
                "title": "Workplace disputes and working conditions"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peopleinwork/workplacepensions",
              "description": {
                "title": "Workplace pensions"
              },
              "type": "product_page"
            }
          ],
          "title": "People in work"
        },
        {
          "uri": "/employmentandlabourmarket/peoplenotinwork",
          "description": {},
          "children": [
            {
              "uri": "/employmentandlabourmarket/peoplenotinwork/economicinactivity",
              "description": {
                "title": "Economic inactivity"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peoplenotinwork/outofworkbenefits",
              "description": {
                "title": "Out of work benefits"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peoplenotinwork/redundancies",
              "description": {
                "title": "Redundancies"
              },
              "type": "product_page"
            },
            {
              "uri": "/employmentandlabourmarket/peoplenotinwork/unemployment",
              "description": {
                "title": "Unemployment"
              },
              "type": "product_page"
            }
          ],
          "title": "People not in work"
        }
      ],
      "title": "Employment and labour market"
    },
    {
      "uri": "/peoplepopulationandcommunity",
      "children": [
        {
          "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages",
          "description": {},
          "children": [
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/adoption",
              "description": {
                "title": "Adoption"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/ageing",
              "description": {
                "title": "Ageing"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/conceptionandfertilityrates",
              "description": {
                "title": "Conception and fertility rates"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/deaths",
              "description": {
                "title": "Deaths"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/divorce",
              "description": {
                "title": "Divorce"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/families",
              "description": {
                "title": "Families"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/lifeexpectancies",
              "description": {
                "title": "Life expectancies"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/livebirths",
              "description": {
                "title": "Live births"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/marriagecohabitationandcivilpartnerships",
              "description": {
                "title": "Marriage, cohabitation and civil partnerships"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/maternities",
              "description": {
                "title": "Maternities"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/birthsdeathsandmarriages/stillbirths",
              "description": {
                "title": "Stillbirths"
              },
              "type": "product_page"
            }
          ],
          "title": "Births, deaths and marriages"
        },
        {
          "uri": "/peoplepopulationandcommunity/crimeandjustice",
          "description": {},
          "title": "Crime and justice"
        },
        {
          "uri": "/peoplepopulationandcommunity/culturalidentity",
          "description": {},
          "children": [
            {
              "uri": "/peoplepopulationandcommunity/culturalidentity/ethnicity",
              "description": {
                "title": "Ethnicity"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/culturalidentity/language",
              "description": {
                "title": "Language"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/culturalidentity/religion",
              "description": {
                "title": "Religion"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/culturalidentity/sexuality",
              "description": {
                "title": "Sexuality"
              },
              "type": "product_page"
            }
          ],
          "title": "Cultural identity"
        },
        {
          "uri": "/peoplepopulationandcommunity/educationandchildcare",
          "description": {},
          "title": "Education and childcare"
        },
        {
          "uri": "/peoplepopulationandcommunity/elections",
          "description": {},
          "children": [
            {
              "uri": "/peoplepopulationandcommunity/elections/electoralregistration",
              "description": {
                "title": "Electoral registration"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/elections/generalelections",
              "description": {
                "title": "General elections"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/elections/localgovernmentelections",
              "description": {
                "title": "Local government elections"
              },
              "type": "product_page"
            }
          ],
          "title": "Elections"
        },
        {
          "uri": "/peoplepopulationandcommunity/healthandsocialcare",
          "description": {},
          "children": [
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/causesofdeath",
              "description": {
                "title": "Causes of death"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/childhealth",
              "description": {
                "title": "Child health"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/conditionsanddiseases",
              "description": {
                "title": "Conditions and diseases"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/disability",
              "description": {
                "title": "Disability"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/drugusealcoholandsmoking",
              "description": {
                "title": "Drug use, alcohol and smoking"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/healthandlifeexpectancies",
              "description": {
                "title": "Health and life expectancies"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/healthandwellbeing",
              "description": {
                "title": "Health and well-being"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/healthcaresystem",
              "description": {
                "title": "Health care system"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/healthinequalities",
              "description": {
                "title": "Health inequalities"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/mentalhealth",
              "description": {
                "title": "Mental health"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/healthandsocialcare/socialcare",
              "description": {
                "title": "Social care"
              },
              "type": "product_page"
            }
          ],
          "title": "Health and social care"
        },
        {
          "uri": "/peoplepopulationandcommunity/householdcharacteristics",
          "description": {},
          "children": [
            {
              "uri": "/peoplepopulationandcommunity/householdcharacteristics/homeinternetandsocialmediausage",
              "description": {
                "title": "Home internet and social media usage"
              },
              "type": "product_page"
            }
          ],
          "title": "Household characteristics"
        },
        {
          "uri": "/peoplepopulationandcommunity/housing",
          "description": {},
          "title": "Housing"
        },
        {
          "uri": "/peoplepopulationandcommunity/leisureandtourism",
          "description": {},
          "title": "Leisure and tourism"
        },
        {
          "uri": "/peoplepopulationandcommunity/personalandhouseholdfinances",
          "description": {},
          "children": [
            {
              "uri": "/peoplepopulationandcommunity/personalandhouseholdfinances/debt",
              "description": {
                "title": "Debt"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/personalandhouseholdfinances/expenditure",
              "description": {
                "title": "Expenditure"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/personalandhouseholdfinances/incomeandwealth",
              "description": {
                "title": "Income and wealth"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/personalandhouseholdfinances/pensionssavingsandinvestments",
              "description": {
                "title": "Pensions, savings and investments"
              },
              "type": "product_page"
            }
          ],
          "title": "Personal and household finances"
        },
        {
          "uri": "/peoplepopulationandcommunity/populationandmigration",
          "description": {},
          "children": [
            {
              "uri": "/peoplepopulationandcommunity/populationandmigration/internationalmigration",
              "description": {
                "title": "International migration"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/populationandmigration/migrationwithintheuk",
              "description": {
                "title": "Migration within the UK"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/populationandmigration/populationestimates",
              "description": {
                "title": "Population estimates"
              },
              "type": "product_page"
            },
            {
              "uri": "/peoplepopulationandcommunity/populationandmigration/populationprojections",
              "description": {
                "title": "Population projections"
              },
              "type": "product_page"
            }
          ],
          "title": "Population and migration"
        },
        {
          "uri": "/peoplepopulationandcommunity/wellbeing",
          "description": {},
          "title": "Well-being"
        }
      ],
      "title": "People, population and community"
    }
  ],
  "uri": "/",
  "type": "homepage",
  "metadata": {
    "title": "Home",
    "description": "thing",
    "keywords": [
      "statistics",
      "economy",
      "census",
      "population",
      "inflation",
      "employment"
    ]
  },
  "data": {
    "publications": [
      {
        "title": "Construction output in Great Britain: Aug 2016",
        "uri": "/businessindustryandtrade/constructionindustry/bulletins/constructionoutputingreatbritain/aug2016",
        "releaseDate": "13 October 2016"
      },
      {
        "title": "Profitability of UK companies: Apr to June 2016",
        "uri": "/economy/nationalaccounts/uksectoraccounts/bulletins/profitabilityofukcompanies/aprtojun2016",
        "releaseDate": "12 October 2016"
      },
      {
        "title": "Overseas travel and tourism, provisional: Apr to June 2016",
        "uri": "/peoplepopulationandcommunity/leisureandtourism/articles/overseastravelandtourismprovisionalresults/apriltojune2016",
        "releaseDate": "12 October 2016"
      },
      {
        "title": "Alternative measures of real households disposable income and the saving ratio: Sept 2016",
        "uri": "/peoplepopulationandcommunity/personalandhouseholdfinances/incomeandwealth/articles/alternativemeasuresofrealhouseholdsdisposableincomeandthesavingratio/sept2016",
        "releaseDate": "11 October 2016"
      },
      {
        "title": "National Accounts: proposed methodological changes to chainlinking for UK publications and international transmissions",
        "uri": "/peoplepopulationandcommunity/personalandhouseholdfinances/incomeandwealth/articles/alternativemeasuresofrealhouseholdsdisposableincomeandthesavingratio/sept2016",
        "releaseDate": "11 October 2016"
      },
      {
        "title": "National Accounts: proposed methodological changes to chainlinking for UK publications and international transmissions",
        "uri": "/economy/grossdomesticproductgdp/articles/nationalaccounts/proposedmethodologicalchangestochainlinkingforukpublicationsandinternationaltransmissions",
        "releaseDate": "11 October 2016"
      }
    ],
    "headlineFigures": [
      {
        "title": "CPI: Consumer Prices Index (% change)",
        "uri": "/economy/inflationandpriceindices/timeseries/d7g7/mm23",
        "releaseDate": "Jun 2016",
        "latestFigure": {
          "figure": "0.5",
          "preUnit": "",
          "unit": "%"
        },
        "sparklineData": [
          {
            "name": "1989 JAN",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "1989 FEB",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "1989 MAR",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "1989 APR",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1989 MAY",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1989 JUN",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "1989 JUL",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "1989 AUG",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "1989 SEP",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "1989 OCT",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1989 NOV",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1989 DEC",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1990 JAN",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "1990 FEB",
            "y": 5.9,
            "stringY": "5.9"
          },
          {
            "name": "1990 MAR",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "1990 APR",
            "y": 6.4,
            "stringY": "6.4"
          },
          {
            "name": "1990 MAY",
            "y": 6.8,
            "stringY": "6.8"
          },
          {
            "name": "1990 JUN",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1990 JUL",
            "y": 6.8,
            "stringY": "6.8"
          },
          {
            "name": "1990 AUG",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "1990 SEP",
            "y": 8.1,
            "stringY": "8.1"
          },
          {
            "name": "1990 OCT",
            "y": 8.1,
            "stringY": "8.1"
          },
          {
            "name": "1990 NOV",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "1990 DEC",
            "y": 7.6,
            "stringY": "7.6"
          },
          {
            "name": "1991 JAN",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1991 FEB",
            "y": 7,
            "stringY": "7.0"
          },
          {
            "name": "1991 MAR",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1991 APR",
            "y": 8.5,
            "stringY": "8.5"
          },
          {
            "name": "1991 MAY",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "1991 JUN",
            "y": 8.4,
            "stringY": "8.4"
          },
          {
            "name": "1991 JUL",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1991 AUG",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "1991 SEP",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1991 OCT",
            "y": 6.8,
            "stringY": "6.8"
          },
          {
            "name": "1991 NOV",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1991 DEC",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "1992 JAN",
            "y": 7,
            "stringY": "7.0"
          },
          {
            "name": "1992 FEB",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1992 MAR",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1992 APR",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "1992 MAY",
            "y": 4.3,
            "stringY": "4.3"
          },
          {
            "name": "1992 JUN",
            "y": 3.8,
            "stringY": "3.8"
          },
          {
            "name": "1992 JUL",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1992 AUG",
            "y": 3.2,
            "stringY": "3.2"
          },
          {
            "name": "1992 SEP",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "1992 OCT",
            "y": 2.9,
            "stringY": "2.9"
          },
          {
            "name": "1992 NOV",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1992 DEC",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1993 JAN",
            "y": 2.2,
            "stringY": "2.2"
          },
          {
            "name": "1993 FEB",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "1993 MAR",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "1993 APR",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "1993 MAY",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "1993 JUN",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "1993 JUL",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1993 AUG",
            "y": 2.9,
            "stringY": "2.9"
          },
          {
            "name": "1993 SEP",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "1993 OCT",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1993 NOV",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "1993 DEC",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "1994 JAN",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "1994 FEB",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "1994 MAR",
            "y": 2.2,
            "stringY": "2.2"
          },
          {
            "name": "1994 APR",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "1994 MAY",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "1994 JUN",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "1994 JUL",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "1994 AUG",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "1994 SEP",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "1994 OCT",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "1994 NOV",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "1994 DEC",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "1995 JAN",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "1995 FEB",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "1995 MAR",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1995 APR",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "1995 MAY",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "1995 JUN",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1995 JUL",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1995 AUG",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1995 SEP",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "1995 OCT",
            "y": 2.9,
            "stringY": "2.9"
          },
          {
            "name": "1995 NOV",
            "y": 2.8,
            "stringY": "2.8"
          },
          {
            "name": "1995 DEC",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "1996 JAN",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "1996 FEB",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "1996 MAR",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1996 APR",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1996 MAY",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "1996 JUN",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "1996 JUL",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "1996 AUG",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "1996 SEP",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "1996 OCT",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "1996 NOV",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "1996 DEC",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "1997 JAN",
            "y": 2.1,
            "stringY": "2.1"
          },
          {
            "name": "1997 FEB",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "1997 MAR",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "1997 APR",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "1997 MAY",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "1997 JUN",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "1997 JUL",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "1997 AUG",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "1997 SEP",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "1997 OCT",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "1997 NOV",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "1997 DEC",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "1998 JAN",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "1998 FEB",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "1998 MAR",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "1998 APR",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "1998 MAY",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "1998 JUN",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "1998 JUL",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "1998 AUG",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "1998 SEP",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "1998 OCT",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "1998 NOV",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "1998 DEC",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "1999 JAN",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "1999 FEB",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "1999 MAR",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "1999 APR",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "1999 MAY",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "1999 JUN",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "1999 JUL",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "1999 AUG",
            "y": 1.2,
            "stringY": "1.2"
          },
          {
            "name": "1999 SEP",
            "y": 1.2,
            "stringY": "1.2"
          },
          {
            "name": "1999 OCT",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "1999 NOV",
            "y": 1.2,
            "stringY": "1.2"
          },
          {
            "name": "1999 DEC",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2000 JAN",
            "y": 0.8,
            "stringY": "0.8"
          },
          {
            "name": "2000 FEB",
            "y": 0.9,
            "stringY": "0.9"
          },
          {
            "name": "2000 MAR",
            "y": 0.6,
            "stringY": "0.6"
          },
          {
            "name": "2000 APR",
            "y": 0.6,
            "stringY": "0.6"
          },
          {
            "name": "2000 MAY",
            "y": 0.5,
            "stringY": "0.5"
          },
          {
            "name": "2000 JUN",
            "y": 0.8,
            "stringY": "0.8"
          },
          {
            "name": "2000 JUL",
            "y": 0.9,
            "stringY": "0.9"
          },
          {
            "name": "2000 AUG",
            "y": 0.6,
            "stringY": "0.6"
          },
          {
            "name": "2000 SEP",
            "y": 1,
            "stringY": "1.0"
          },
          {
            "name": "2000 OCT",
            "y": 1,
            "stringY": "1.0"
          },
          {
            "name": "2000 NOV",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2000 DEC",
            "y": 0.8,
            "stringY": "0.8"
          },
          {
            "name": "2001 JAN",
            "y": 0.9,
            "stringY": "0.9"
          },
          {
            "name": "2001 FEB",
            "y": 0.8,
            "stringY": "0.8"
          },
          {
            "name": "2001 MAR",
            "y": 0.9,
            "stringY": "0.9"
          },
          {
            "name": "2001 APR",
            "y": 1.2,
            "stringY": "1.2"
          },
          {
            "name": "2001 MAY",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "2001 JUN",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "2001 JUL",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2001 AUG",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "2001 SEP",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2001 OCT",
            "y": 1.2,
            "stringY": "1.2"
          },
          {
            "name": "2001 NOV",
            "y": 0.8,
            "stringY": "0.8"
          },
          {
            "name": "2001 DEC",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2002 JAN",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "2002 FEB",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2002 MAR",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2002 APR",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2002 MAY",
            "y": 0.8,
            "stringY": "0.8"
          },
          {
            "name": "2002 JUN",
            "y": 0.6,
            "stringY": "0.6"
          },
          {
            "name": "2002 JUL",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2002 AUG",
            "y": 1,
            "stringY": "1.0"
          },
          {
            "name": "2002 SEP",
            "y": 1,
            "stringY": "1.0"
          },
          {
            "name": "2002 OCT",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2002 NOV",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2002 DEC",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "2003 JAN",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2003 FEB",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "2003 MAR",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2003 APR",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2003 MAY",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2003 JUN",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2003 JUL",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2003 AUG",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2003 SEP",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2003 OCT",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2003 NOV",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2003 DEC",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2004 JAN",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2004 FEB",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2004 MAR",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2004 APR",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2004 MAY",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2004 JUN",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "2004 JUL",
            "y": 1.4,
            "stringY": "1.4"
          },
          {
            "name": "2004 AUG",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2004 SEP",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2004 OCT",
            "y": 1.2,
            "stringY": "1.2"
          },
          {
            "name": "2004 NOV",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2004 DEC",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "2005 JAN",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "2005 FEB",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "2005 MAR",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2005 APR",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2005 MAY",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2005 JUN",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "2005 JUL",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "2005 AUG",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "2005 SEP",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "2005 OCT",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "2005 NOV",
            "y": 2.1,
            "stringY": "2.1"
          },
          {
            "name": "2005 DEC",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2006 JAN",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2006 FEB",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "2006 MAR",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "2006 APR",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "2006 MAY",
            "y": 2.2,
            "stringY": "2.2"
          },
          {
            "name": "2006 JUN",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "2006 JUL",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "2006 AUG",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "2006 SEP",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "2006 OCT",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "2006 NOV",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2006 DEC",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "2007 JAN",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2007 FEB",
            "y": 2.8,
            "stringY": "2.8"
          },
          {
            "name": "2007 MAR",
            "y": 3.1,
            "stringY": "3.1"
          },
          {
            "name": "2007 APR",
            "y": 2.8,
            "stringY": "2.8"
          },
          {
            "name": "2007 MAY",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "2007 JUN",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "2007 JUL",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2007 AUG",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "2007 SEP",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "2007 OCT",
            "y": 2.1,
            "stringY": "2.1"
          },
          {
            "name": "2007 NOV",
            "y": 2.1,
            "stringY": "2.1"
          },
          {
            "name": "2007 DEC",
            "y": 2.1,
            "stringY": "2.1"
          },
          {
            "name": "2008 JAN",
            "y": 2.2,
            "stringY": "2.2"
          },
          {
            "name": "2008 FEB",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "2008 MAR",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "2008 APR",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "2008 MAY",
            "y": 3.3,
            "stringY": "3.3"
          },
          {
            "name": "2008 JUN",
            "y": 3.8,
            "stringY": "3.8"
          },
          {
            "name": "2008 JUL",
            "y": 4.4,
            "stringY": "4.4"
          },
          {
            "name": "2008 AUG",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2008 SEP",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2008 OCT",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "2008 NOV",
            "y": 4.1,
            "stringY": "4.1"
          },
          {
            "name": "2008 DEC",
            "y": 3.1,
            "stringY": "3.1"
          },
          {
            "name": "2009 JAN",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "2009 FEB",
            "y": 3.2,
            "stringY": "3.2"
          },
          {
            "name": "2009 MAR",
            "y": 2.9,
            "stringY": "2.9"
          },
          {
            "name": "2009 APR",
            "y": 2.3,
            "stringY": "2.3"
          },
          {
            "name": "2009 MAY",
            "y": 2.2,
            "stringY": "2.2"
          },
          {
            "name": "2009 JUN",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "2009 JUL",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "2009 AUG",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "2009 SEP",
            "y": 1.1,
            "stringY": "1.1"
          },
          {
            "name": "2009 OCT",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2009 NOV",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2009 DEC",
            "y": 2.9,
            "stringY": "2.9"
          },
          {
            "name": "2010 JAN",
            "y": 3.5,
            "stringY": "3.5"
          },
          {
            "name": "2010 FEB",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "2010 MAR",
            "y": 3.4,
            "stringY": "3.4"
          },
          {
            "name": "2010 APR",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "2010 MAY",
            "y": 3.4,
            "stringY": "3.4"
          },
          {
            "name": "2010 JUN",
            "y": 3.2,
            "stringY": "3.2"
          },
          {
            "name": "2010 JUL",
            "y": 3.1,
            "stringY": "3.1"
          },
          {
            "name": "2010 AUG",
            "y": 3.1,
            "stringY": "3.1"
          },
          {
            "name": "2010 SEP",
            "y": 3.1,
            "stringY": "3.1"
          },
          {
            "name": "2010 OCT",
            "y": 3.2,
            "stringY": "3.2"
          },
          {
            "name": "2010 NOV",
            "y": 3.3,
            "stringY": "3.3"
          },
          {
            "name": "2010 DEC",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "2011 JAN",
            "y": 4,
            "stringY": "4.0"
          },
          {
            "name": "2011 FEB",
            "y": 4.4,
            "stringY": "4.4"
          },
          {
            "name": "2011 MAR",
            "y": 4,
            "stringY": "4.0"
          },
          {
            "name": "2011 APR",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "2011 MAY",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "2011 JUN",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "2011 JUL",
            "y": 4.4,
            "stringY": "4.4"
          },
          {
            "name": "2011 AUG",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "2011 SEP",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2011 OCT",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2011 NOV",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2011 DEC",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "2012 JAN",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "2012 FEB",
            "y": 3.4,
            "stringY": "3.4"
          },
          {
            "name": "2012 MAR",
            "y": 3.5,
            "stringY": "3.5"
          },
          {
            "name": "2012 APR",
            "y": 3,
            "stringY": "3.0"
          },
          {
            "name": "2012 MAY",
            "y": 2.8,
            "stringY": "2.8"
          },
          {
            "name": "2012 JUN",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "2012 JUL",
            "y": 2.6,
            "stringY": "2.6"
          },
          {
            "name": "2012 AUG",
            "y": 2.5,
            "stringY": "2.5"
          },
          {
            "name": "2012 SEP",
            "y": 2.2,
            "stringY": "2.2"
          },
          {
            "name": "2012 OCT",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2012 NOV",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2012 DEC",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2013 JAN",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2013 FEB",
            "y": 2.8,
            "stringY": "2.8"
          },
          {
            "name": "2013 MAR",
            "y": 2.8,
            "stringY": "2.8"
          },
          {
            "name": "2013 APR",
            "y": 2.4,
            "stringY": "2.4"
          },
          {
            "name": "2013 MAY",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2013 JUN",
            "y": 2.9,
            "stringY": "2.9"
          },
          {
            "name": "2013 JUL",
            "y": 2.8,
            "stringY": "2.8"
          },
          {
            "name": "2013 AUG",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2013 SEP",
            "y": 2.7,
            "stringY": "2.7"
          },
          {
            "name": "2013 OCT",
            "y": 2.2,
            "stringY": "2.2"
          },
          {
            "name": "2013 NOV",
            "y": 2.1,
            "stringY": "2.1"
          },
          {
            "name": "2013 DEC",
            "y": 2,
            "stringY": "2.0"
          },
          {
            "name": "2014 JAN",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2014 FEB",
            "y": 1.7,
            "stringY": "1.7"
          },
          {
            "name": "2014 MAR",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "2014 APR",
            "y": 1.8,
            "stringY": "1.8"
          },
          {
            "name": "2014 MAY",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2014 JUN",
            "y": 1.9,
            "stringY": "1.9"
          },
          {
            "name": "2014 JUL",
            "y": 1.6,
            "stringY": "1.6"
          },
          {
            "name": "2014 AUG",
            "y": 1.5,
            "stringY": "1.5"
          },
          {
            "name": "2014 SEP",
            "y": 1.2,
            "stringY": "1.2"
          },
          {
            "name": "2014 OCT",
            "y": 1.3,
            "stringY": "1.3"
          },
          {
            "name": "2014 NOV",
            "y": 1,
            "stringY": "1.0"
          },
          {
            "name": "2014 DEC",
            "y": 0.5,
            "stringY": "0.5"
          },
          {
            "name": "2015 JAN",
            "y": 0.3,
            "stringY": "0.3"
          },
          {
            "name": "2015 FEB",
            "y": 0,
            "stringY": "0.0"
          },
          {
            "name": "2015 MAR",
            "y": 0,
            "stringY": "0.0"
          },
          {
            "name": "2015 APR",
            "y": -0.1,
            "stringY": "-0.1"
          },
          {
            "name": "2015 MAY",
            "y": 0.1,
            "stringY": "0.1"
          },
          {
            "name": "2015 JUN",
            "y": 0,
            "stringY": "0.0"
          },
          {
            "name": "2015 JUL",
            "y": 0.1,
            "stringY": "0.1"
          },
          {
            "name": "2015 AUG",
            "y": 0,
            "stringY": "0.0"
          },
          {
            "name": "2015 SEP",
            "y": -0.1,
            "stringY": "-0.1"
          },
          {
            "name": "2015 OCT",
            "y": -0.1,
            "stringY": "-0.1"
          },
          {
            "name": "2015 NOV",
            "y": 0.1,
            "stringY": "0.1"
          },
          {
            "name": "2015 DEC",
            "y": 0.2,
            "stringY": "0.2"
          },
          {
            "name": "2016 JAN",
            "y": 0.3,
            "stringY": "0.3"
          },
          {
            "name": "2016 FEB",
            "y": 0.3,
            "stringY": "0.3"
          },
          {
            "name": "2016 MAR",
            "y": 0.5,
            "stringY": "0.5"
          },
          {
            "name": "2016 APR",
            "y": 0.3,
            "stringY": "0.3"
          },
          {
            "name": "2016 MAY",
            "y": 0.3,
            "stringY": "0.3"
          },
          {
            "name": "2016 JUN",
            "y": 0.5,
            "stringY": "0.5"
          },
          {
            "name": "2016 JUL",
            "y": 0.6,
            "stringY": "0.6"
          },
          {
            "name": "2016 AUG",
            "y": 0.6,
            "stringY": "0.6"
          }
        ]
      },
      {
        "title": "Gross Domestic Product: quarter on quarter growth: CVM SA %",
        "uri": "/employmentandlabourmarket/peopleinwork/employmentandemployeetypes/timeseries/lf24/lms",
        "releaseDate": "Q2 2016",
        "latestFigure": {
          "figure": "0.6",
          "preUnit": "",
          "unit": "%"
        },
        "sparklineData": [
          {
            "name": "1971 JAN-MAR",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "1971 FEB-APR",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "1971 MAR-MAY",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1971 APR-JUN",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1971 MAY-JUL",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1971 JUN-AUG",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1971 JUL-SEP",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1971 AUG-OCT",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1971 SEP-NOV",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1971 OCT-DEC",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1971 NOV-JAN",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1971 DEC-FEB",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1972 JAN-MAR",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1972 FEB-APR",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1972 MAR-MAY",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1972 APR-JUN",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1972 MAY-JUL",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1972 JUN-AUG",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1972 JUL-SEP",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1972 AUG-OCT",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1972 SEP-NOV",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "1972 OCT-DEC",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "1972 NOV-JAN",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "1972 DEC-FEB",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "1973 JAN-MAR",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1973 FEB-APR",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1973 MAR-MAY",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1973 APR-JUN",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1973 MAY-JUL",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1973 JUN-AUG",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1973 JUL-SEP",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1973 AUG-OCT",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1973 SEP-NOV",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1973 OCT-DEC",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1973 NOV-JAN",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1973 DEC-FEB",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1974 JAN-MAR",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1974 FEB-APR",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1974 MAR-MAY",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1974 APR-JUN",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1974 MAY-JUL",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1974 JUN-AUG",
            "y": 73.1,
            "stringY": "73.1"
          },
          {
            "name": "1974 JUL-SEP",
            "y": 73.1,
            "stringY": "73.1"
          },
          {
            "name": "1974 AUG-OCT",
            "y": 73.1,
            "stringY": "73.1"
          },
          {
            "name": "1974 SEP-NOV",
            "y": 73.1,
            "stringY": "73.1"
          },
          {
            "name": "1974 OCT-DEC",
            "y": 73.1,
            "stringY": "73.1"
          },
          {
            "name": "1974 NOV-JAN",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "1974 DEC-FEB",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1975 JAN-MAR",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "1975 FEB-APR",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1975 MAR-MAY",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "1975 APR-JUN",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "1975 MAY-JUL",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "1975 JUN-AUG",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "1975 JUL-SEP",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "1975 AUG-OCT",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "1975 SEP-NOV",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "1975 OCT-DEC",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "1975 NOV-JAN",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "1975 DEC-FEB",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "1976 JAN-MAR",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "1976 FEB-APR",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1976 MAR-MAY",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1976 APR-JUN",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1976 MAY-JUL",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1976 JUN-AUG",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1976 JUL-SEP",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1976 AUG-OCT",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1976 SEP-NOV",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1976 OCT-DEC",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1976 NOV-JAN",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1976 DEC-FEB",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1977 JAN-MAR",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1977 FEB-APR",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1977 MAR-MAY",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1977 APR-JUN",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1977 MAY-JUL",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1977 JUN-AUG",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1977 JUL-SEP",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1977 AUG-OCT",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1977 SEP-NOV",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1977 OCT-DEC",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1977 NOV-JAN",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1977 DEC-FEB",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 JAN-MAR",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 FEB-APR",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 MAR-MAY",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 APR-JUN",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 MAY-JUL",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 JUN-AUG",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 JUL-SEP",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1978 AUG-OCT",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1978 SEP-NOV",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1978 OCT-DEC",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1978 NOV-JAN",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1978 DEC-FEB",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1979 JAN-MAR",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1979 FEB-APR",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1979 MAR-MAY",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1979 APR-JUN",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1979 MAY-JUL",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1979 JUN-AUG",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1979 JUL-SEP",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1979 AUG-OCT",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1979 SEP-NOV",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1979 OCT-DEC",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1979 NOV-JAN",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1979 DEC-FEB",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1980 JAN-MAR",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1980 FEB-APR",
            "y": 71.4,
            "stringY": "71.4"
          },
          {
            "name": "1980 MAR-MAY",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "1980 APR-JUN",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "1980 MAY-JUL",
            "y": 71.1,
            "stringY": "71.1"
          },
          {
            "name": "1980 JUN-AUG",
            "y": 70.9,
            "stringY": "70.9"
          },
          {
            "name": "1980 JUL-SEP",
            "y": 70.7,
            "stringY": "70.7"
          },
          {
            "name": "1980 AUG-OCT",
            "y": 70.4,
            "stringY": "70.4"
          },
          {
            "name": "1980 SEP-NOV",
            "y": 70.2,
            "stringY": "70.2"
          },
          {
            "name": "1980 OCT-DEC",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1980 NOV-JAN",
            "y": 69.7,
            "stringY": "69.7"
          },
          {
            "name": "1980 DEC-FEB",
            "y": 69.5,
            "stringY": "69.5"
          },
          {
            "name": "1981 JAN-MAR",
            "y": 69.3,
            "stringY": "69.3"
          },
          {
            "name": "1981 FEB-APR",
            "y": 69.1,
            "stringY": "69.1"
          },
          {
            "name": "1981 MAR-MAY",
            "y": 68.9,
            "stringY": "68.9"
          },
          {
            "name": "1981 APR-JUN",
            "y": 68.7,
            "stringY": "68.7"
          },
          {
            "name": "1981 MAY-JUL",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1981 JUN-AUG",
            "y": 68.4,
            "stringY": "68.4"
          },
          {
            "name": "1981 JUL-SEP",
            "y": 68.2,
            "stringY": "68.2"
          },
          {
            "name": "1981 AUG-OCT",
            "y": 68.1,
            "stringY": "68.1"
          },
          {
            "name": "1981 SEP-NOV",
            "y": 67.9,
            "stringY": "67.9"
          },
          {
            "name": "1981 OCT-DEC",
            "y": 67.7,
            "stringY": "67.7"
          },
          {
            "name": "1981 NOV-JAN",
            "y": 67.6,
            "stringY": "67.6"
          },
          {
            "name": "1981 DEC-FEB",
            "y": 67.5,
            "stringY": "67.5"
          },
          {
            "name": "1982 JAN-MAR",
            "y": 67.4,
            "stringY": "67.4"
          },
          {
            "name": "1982 FEB-APR",
            "y": 67.3,
            "stringY": "67.3"
          },
          {
            "name": "1982 MAR-MAY",
            "y": 67.2,
            "stringY": "67.2"
          },
          {
            "name": "1982 APR-JUN",
            "y": 67.1,
            "stringY": "67.1"
          },
          {
            "name": "1982 MAY-JUL",
            "y": 66.9,
            "stringY": "66.9"
          },
          {
            "name": "1982 JUN-AUG",
            "y": 66.8,
            "stringY": "66.8"
          },
          {
            "name": "1982 JUL-SEP",
            "y": 66.6,
            "stringY": "66.6"
          },
          {
            "name": "1982 AUG-OCT",
            "y": 66.5,
            "stringY": "66.5"
          },
          {
            "name": "1982 SEP-NOV",
            "y": 66.3,
            "stringY": "66.3"
          },
          {
            "name": "1982 OCT-DEC",
            "y": 66.1,
            "stringY": "66.1"
          },
          {
            "name": "1982 NOV-JAN",
            "y": 65.9,
            "stringY": "65.9"
          },
          {
            "name": "1982 DEC-FEB",
            "y": 65.8,
            "stringY": "65.8"
          },
          {
            "name": "1983 JAN-MAR",
            "y": 65.7,
            "stringY": "65.7"
          },
          {
            "name": "1983 FEB-APR",
            "y": 65.6,
            "stringY": "65.6"
          },
          {
            "name": "1983 MAR-MAY",
            "y": 65.6,
            "stringY": "65.6"
          },
          {
            "name": "1983 APR-JUN",
            "y": 65.6,
            "stringY": "65.6"
          },
          {
            "name": "1983 MAY-JUL",
            "y": 65.6,
            "stringY": "65.6"
          },
          {
            "name": "1983 JUN-AUG",
            "y": 65.7,
            "stringY": "65.7"
          },
          {
            "name": "1983 JUL-SEP",
            "y": 65.9,
            "stringY": "65.9"
          },
          {
            "name": "1983 AUG-OCT",
            "y": 66,
            "stringY": "66.0"
          },
          {
            "name": "1983 SEP-NOV",
            "y": 66.2,
            "stringY": "66.2"
          },
          {
            "name": "1983 OCT-DEC",
            "y": 66.3,
            "stringY": "66.3"
          },
          {
            "name": "1983 NOV-JAN",
            "y": 66.4,
            "stringY": "66.4"
          },
          {
            "name": "1983 DEC-FEB",
            "y": 66.4,
            "stringY": "66.4"
          },
          {
            "name": "1984 JAN-MAR",
            "y": 66.5,
            "stringY": "66.5"
          },
          {
            "name": "1984 FEB-APR",
            "y": 66.6,
            "stringY": "66.6"
          },
          {
            "name": "1984 MAR-MAY",
            "y": 66.6,
            "stringY": "66.6"
          },
          {
            "name": "1984 APR-JUN",
            "y": 66.7,
            "stringY": "66.7"
          },
          {
            "name": "1984 MAY-JUL",
            "y": 66.7,
            "stringY": "66.7"
          },
          {
            "name": "1984 JUN-AUG",
            "y": 66.8,
            "stringY": "66.8"
          },
          {
            "name": "1984 JUL-SEP",
            "y": 66.8,
            "stringY": "66.8"
          },
          {
            "name": "1984 AUG-OCT",
            "y": 66.9,
            "stringY": "66.9"
          },
          {
            "name": "1984 SEP-NOV",
            "y": 67,
            "stringY": "67.0"
          },
          {
            "name": "1984 OCT-DEC",
            "y": 67.1,
            "stringY": "67.1"
          },
          {
            "name": "1984 NOV-JAN",
            "y": 67.1,
            "stringY": "67.1"
          },
          {
            "name": "1984 DEC-FEB",
            "y": 67.2,
            "stringY": "67.2"
          },
          {
            "name": "1985 JAN-MAR",
            "y": 67.3,
            "stringY": "67.3"
          },
          {
            "name": "1985 FEB-APR",
            "y": 67.3,
            "stringY": "67.3"
          },
          {
            "name": "1985 MAR-MAY",
            "y": 67.3,
            "stringY": "67.3"
          },
          {
            "name": "1985 APR-JUN",
            "y": 67.4,
            "stringY": "67.4"
          },
          {
            "name": "1985 MAY-JUL",
            "y": 67.4,
            "stringY": "67.4"
          },
          {
            "name": "1985 JUN-AUG",
            "y": 67.4,
            "stringY": "67.4"
          },
          {
            "name": "1985 JUL-SEP",
            "y": 67.5,
            "stringY": "67.5"
          },
          {
            "name": "1985 AUG-OCT",
            "y": 67.5,
            "stringY": "67.5"
          },
          {
            "name": "1985 SEP-NOV",
            "y": 67.5,
            "stringY": "67.5"
          },
          {
            "name": "1985 OCT-DEC",
            "y": 67.5,
            "stringY": "67.5"
          },
          {
            "name": "1985 NOV-JAN",
            "y": 67.5,
            "stringY": "67.5"
          },
          {
            "name": "1985 DEC-FEB",
            "y": 67.6,
            "stringY": "67.6"
          },
          {
            "name": "1986 JAN-MAR",
            "y": 67.6,
            "stringY": "67.6"
          },
          {
            "name": "1986 FEB-APR",
            "y": 67.6,
            "stringY": "67.6"
          },
          {
            "name": "1986 MAR-MAY",
            "y": 67.6,
            "stringY": "67.6"
          },
          {
            "name": "1986 APR-JUN",
            "y": 67.6,
            "stringY": "67.6"
          },
          {
            "name": "1986 MAY-JUL",
            "y": 67.6,
            "stringY": "67.6"
          },
          {
            "name": "1986 JUN-AUG",
            "y": 67.7,
            "stringY": "67.7"
          },
          {
            "name": "1986 JUL-SEP",
            "y": 67.8,
            "stringY": "67.8"
          },
          {
            "name": "1986 AUG-OCT",
            "y": 67.8,
            "stringY": "67.8"
          },
          {
            "name": "1986 SEP-NOV",
            "y": 67.9,
            "stringY": "67.9"
          },
          {
            "name": "1986 OCT-DEC",
            "y": 67.9,
            "stringY": "67.9"
          },
          {
            "name": "1986 NOV-JAN",
            "y": 68,
            "stringY": "68.0"
          },
          {
            "name": "1986 DEC-FEB",
            "y": 68,
            "stringY": "68.0"
          },
          {
            "name": "1987 JAN-MAR",
            "y": 68.1,
            "stringY": "68.1"
          },
          {
            "name": "1987 FEB-APR",
            "y": 68.2,
            "stringY": "68.2"
          },
          {
            "name": "1987 MAR-MAY",
            "y": 68.3,
            "stringY": "68.3"
          },
          {
            "name": "1987 APR-JUN",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1987 MAY-JUL",
            "y": 68.7,
            "stringY": "68.7"
          },
          {
            "name": "1987 JUN-AUG",
            "y": 68.9,
            "stringY": "68.9"
          },
          {
            "name": "1987 JUL-SEP",
            "y": 69.1,
            "stringY": "69.1"
          },
          {
            "name": "1987 AUG-OCT",
            "y": 69.3,
            "stringY": "69.3"
          },
          {
            "name": "1987 SEP-NOV",
            "y": 69.5,
            "stringY": "69.5"
          },
          {
            "name": "1987 OCT-DEC",
            "y": 69.7,
            "stringY": "69.7"
          },
          {
            "name": "1987 NOV-JAN",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1987 DEC-FEB",
            "y": 70,
            "stringY": "70.0"
          },
          {
            "name": "1988 JAN-MAR",
            "y": 70.2,
            "stringY": "70.2"
          },
          {
            "name": "1988 FEB-APR",
            "y": 70.4,
            "stringY": "70.4"
          },
          {
            "name": "1988 MAR-MAY",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "1988 APR-JUN",
            "y": 70.7,
            "stringY": "70.7"
          },
          {
            "name": "1988 MAY-JUL",
            "y": 70.8,
            "stringY": "70.8"
          },
          {
            "name": "1988 JUN-AUG",
            "y": 70.9,
            "stringY": "70.9"
          },
          {
            "name": "1988 JUL-SEP",
            "y": 71.1,
            "stringY": "71.1"
          },
          {
            "name": "1988 AUG-OCT",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "1988 SEP-NOV",
            "y": 71.4,
            "stringY": "71.4"
          },
          {
            "name": "1988 OCT-DEC",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1988 NOV-JAN",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1988 DEC-FEB",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1989 JAN-MAR",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "1989 FEB-APR",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "1989 MAR-MAY",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "1989 APR-JUN",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "1989 MAY-JUL",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "1989 JUN-AUG",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "1989 JUL-SEP",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "1989 AUG-OCT",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "1989 SEP-NOV",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "1989 OCT-DEC",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1989 NOV-JAN",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1989 DEC-FEB",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1990 JAN-MAR",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1990 FEB-APR",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1990 MAR-MAY",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1990 APR-JUN",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1990 MAY-JUL",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1990 JUN-AUG",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "1990 JUL-SEP",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "1990 AUG-OCT",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "1990 SEP-NOV",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "1990 OCT-DEC",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "1990 NOV-JAN",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "1990 DEC-FEB",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1991 JAN-MAR",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "1991 FEB-APR",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "1991 MAR-MAY",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "1991 APR-JUN",
            "y": 71,
            "stringY": "71.0"
          },
          {
            "name": "1991 MAY-JUL",
            "y": 70.7,
            "stringY": "70.7"
          },
          {
            "name": "1991 JUN-AUG",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "1991 JUL-SEP",
            "y": 70.3,
            "stringY": "70.3"
          },
          {
            "name": "1991 AUG-OCT",
            "y": 70.1,
            "stringY": "70.1"
          },
          {
            "name": "1991 SEP-NOV",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1991 OCT-DEC",
            "y": 69.8,
            "stringY": "69.8"
          },
          {
            "name": "1991 NOV-JAN",
            "y": 69.6,
            "stringY": "69.6"
          },
          {
            "name": "1991 DEC-FEB",
            "y": 69.5,
            "stringY": "69.5"
          },
          {
            "name": "1992 JAN-MAR",
            "y": 69.4,
            "stringY": "69.4"
          },
          {
            "name": "1992 FEB-APR",
            "y": 69.3,
            "stringY": "69.3"
          },
          {
            "name": "1992 MAR-MAY",
            "y": 69.2,
            "stringY": "69.2"
          },
          {
            "name": "1992 APR-JUN",
            "y": 69.1,
            "stringY": "69.1"
          },
          {
            "name": "1992 MAY-JUL",
            "y": 69,
            "stringY": "69.0"
          },
          {
            "name": "1992 JUN-AUG",
            "y": 68.9,
            "stringY": "68.9"
          },
          {
            "name": "1992 JUL-SEP",
            "y": 68.9,
            "stringY": "68.9"
          },
          {
            "name": "1992 AUG-OCT",
            "y": 68.8,
            "stringY": "68.8"
          },
          {
            "name": "1992 SEP-NOV",
            "y": 68.7,
            "stringY": "68.7"
          },
          {
            "name": "1992 OCT-DEC",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1992 NOV-JAN",
            "y": 68.4,
            "stringY": "68.4"
          },
          {
            "name": "1992 DEC-FEB",
            "y": 68.3,
            "stringY": "68.3"
          },
          {
            "name": "1993 JAN-MAR",
            "y": 68.4,
            "stringY": "68.4"
          },
          {
            "name": "1993 FEB-APR",
            "y": 68.3,
            "stringY": "68.3"
          },
          {
            "name": "1993 MAR-MAY",
            "y": 68.4,
            "stringY": "68.4"
          },
          {
            "name": "1993 APR-JUN",
            "y": 68.4,
            "stringY": "68.4"
          },
          {
            "name": "1993 MAY-JUL",
            "y": 68.4,
            "stringY": "68.4"
          },
          {
            "name": "1993 JUN-AUG",
            "y": 68.4,
            "stringY": "68.4"
          },
          {
            "name": "1993 JUL-SEP",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1993 AUG-OCT",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1993 SEP-NOV",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1993 OCT-DEC",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1993 NOV-JAN",
            "y": 68.5,
            "stringY": "68.5"
          },
          {
            "name": "1993 DEC-FEB",
            "y": 68.6,
            "stringY": "68.6"
          },
          {
            "name": "1994 JAN-MAR",
            "y": 68.7,
            "stringY": "68.7"
          },
          {
            "name": "1994 FEB-APR",
            "y": 68.8,
            "stringY": "68.8"
          },
          {
            "name": "1994 MAR-MAY",
            "y": 68.8,
            "stringY": "68.8"
          },
          {
            "name": "1994 APR-JUN",
            "y": 68.8,
            "stringY": "68.8"
          },
          {
            "name": "1994 MAY-JUL",
            "y": 69,
            "stringY": "69.0"
          },
          {
            "name": "1994 JUN-AUG",
            "y": 69,
            "stringY": "69.0"
          },
          {
            "name": "1994 JUL-SEP",
            "y": 69,
            "stringY": "69.0"
          },
          {
            "name": "1994 AUG-OCT",
            "y": 69,
            "stringY": "69.0"
          },
          {
            "name": "1994 SEP-NOV",
            "y": 69.1,
            "stringY": "69.1"
          },
          {
            "name": "1994 OCT-DEC",
            "y": 69.1,
            "stringY": "69.1"
          },
          {
            "name": "1994 NOV-JAN",
            "y": 69,
            "stringY": "69.0"
          },
          {
            "name": "1994 DEC-FEB",
            "y": 69.1,
            "stringY": "69.1"
          },
          {
            "name": "1995 JAN-MAR",
            "y": 69.2,
            "stringY": "69.2"
          },
          {
            "name": "1995 FEB-APR",
            "y": 69.3,
            "stringY": "69.3"
          },
          {
            "name": "1995 MAR-MAY",
            "y": 69.4,
            "stringY": "69.4"
          },
          {
            "name": "1995 APR-JUN",
            "y": 69.4,
            "stringY": "69.4"
          },
          {
            "name": "1995 MAY-JUL",
            "y": 69.5,
            "stringY": "69.5"
          },
          {
            "name": "1995 JUN-AUG",
            "y": 69.6,
            "stringY": "69.6"
          },
          {
            "name": "1995 JUL-SEP",
            "y": 69.6,
            "stringY": "69.6"
          },
          {
            "name": "1995 AUG-OCT",
            "y": 69.7,
            "stringY": "69.7"
          },
          {
            "name": "1995 SEP-NOV",
            "y": 69.7,
            "stringY": "69.7"
          },
          {
            "name": "1995 OCT-DEC",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1995 NOV-JAN",
            "y": 70,
            "stringY": "70.0"
          },
          {
            "name": "1995 DEC-FEB",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1996 JAN-MAR",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1996 FEB-APR",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1996 MAR-MAY",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1996 APR-JUN",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1996 MAY-JUL",
            "y": 69.9,
            "stringY": "69.9"
          },
          {
            "name": "1996 JUN-AUG",
            "y": 70,
            "stringY": "70.0"
          },
          {
            "name": "1996 JUL-SEP",
            "y": 70,
            "stringY": "70.0"
          },
          {
            "name": "1996 AUG-OCT",
            "y": 70.1,
            "stringY": "70.1"
          },
          {
            "name": "1996 SEP-NOV",
            "y": 70.3,
            "stringY": "70.3"
          },
          {
            "name": "1996 OCT-DEC",
            "y": 70.3,
            "stringY": "70.3"
          },
          {
            "name": "1996 NOV-JAN",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "1996 DEC-FEB",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "1997 JAN-MAR",
            "y": 70.8,
            "stringY": "70.8"
          },
          {
            "name": "1997 FEB-APR",
            "y": 70.9,
            "stringY": "70.9"
          },
          {
            "name": "1997 MAR-MAY",
            "y": 70.9,
            "stringY": "70.9"
          },
          {
            "name": "1997 APR-JUN",
            "y": 71,
            "stringY": "71.0"
          },
          {
            "name": "1997 MAY-JUL",
            "y": 70.9,
            "stringY": "70.9"
          },
          {
            "name": "1997 JUN-AUG",
            "y": 70.9,
            "stringY": "70.9"
          },
          {
            "name": "1997 JUL-SEP",
            "y": 71.1,
            "stringY": "71.1"
          },
          {
            "name": "1997 AUG-OCT",
            "y": 71.1,
            "stringY": "71.1"
          },
          {
            "name": "1997 SEP-NOV",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "1997 OCT-DEC",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "1997 NOV-JAN",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "1997 DEC-FEB",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "1998 JAN-MAR",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "1998 FEB-APR",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "1998 MAR-MAY",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "1998 APR-JUN",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "1998 MAY-JUL",
            "y": 71.4,
            "stringY": "71.4"
          },
          {
            "name": "1998 JUN-AUG",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1998 JUL-SEP",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1998 AUG-OCT",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "1998 SEP-NOV",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1998 OCT-DEC",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1998 NOV-JAN",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1998 DEC-FEB",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1999 JAN-MAR",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1999 FEB-APR",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1999 MAR-MAY",
            "y": 71.8,
            "stringY": "71.8"
          },
          {
            "name": "1999 APR-JUN",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1999 MAY-JUL",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "1999 JUN-AUG",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "1999 JUL-SEP",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "1999 AUG-OCT",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "1999 SEP-NOV",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "1999 OCT-DEC",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "1999 NOV-JAN",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "1999 DEC-FEB",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "2000 JAN-MAR",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "2000 FEB-APR",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "2000 MAR-MAY",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "2000 APR-JUN",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2000 MAY-JUL",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2000 JUN-AUG",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2000 JUL-SEP",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2000 AUG-OCT",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2000 SEP-NOV",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "2000 OCT-DEC",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2000 NOV-JAN",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2000 DEC-FEB",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2001 JAN-MAR",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2001 FEB-APR",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2001 MAR-MAY",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2001 APR-JUN",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2001 MAY-JUL",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2001 JUN-AUG",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2001 JUL-SEP",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2001 AUG-OCT",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2001 SEP-NOV",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2001 OCT-DEC",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2001 NOV-JAN",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2001 DEC-FEB",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2002 JAN-MAR",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2002 FEB-APR",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2002 MAR-MAY",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2002 APR-JUN",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2002 MAY-JUL",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2002 JUN-AUG",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2002 JUL-SEP",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2002 AUG-OCT",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2002 SEP-NOV",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2002 OCT-DEC",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2002 NOV-JAN",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2002 DEC-FEB",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2003 JAN-MAR",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2003 FEB-APR",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2003 MAR-MAY",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2003 APR-JUN",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2003 MAY-JUL",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2003 JUN-AUG",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2003 JUL-SEP",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2003 AUG-OCT",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2003 SEP-NOV",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2003 OCT-DEC",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2003 NOV-JAN",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2003 DEC-FEB",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2004 JAN-MAR",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2004 FEB-APR",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2004 MAR-MAY",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2004 APR-JUN",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2004 MAY-JUL",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2004 JUN-AUG",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2004 JUL-SEP",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2004 AUG-OCT",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2004 SEP-NOV",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2004 OCT-DEC",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2004 NOV-JAN",
            "y": 73.1,
            "stringY": "73.1"
          },
          {
            "name": "2004 DEC-FEB",
            "y": 73.2,
            "stringY": "73.2"
          },
          {
            "name": "2005 JAN-MAR",
            "y": 73.1,
            "stringY": "73.1"
          },
          {
            "name": "2005 FEB-APR",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2005 MAR-MAY",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2005 APR-JUN",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2005 MAY-JUL",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2005 JUN-AUG",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2005 JUL-SEP",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2005 AUG-OCT",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2005 SEP-NOV",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2005 OCT-DEC",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2005 NOV-JAN",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2005 DEC-FEB",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2006 JAN-MAR",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2006 FEB-APR",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2006 MAR-MAY",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2006 APR-JUN",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2006 MAY-JUL",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2006 JUN-AUG",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2006 JUL-SEP",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2006 AUG-OCT",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2006 SEP-NOV",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2006 OCT-DEC",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2006 NOV-JAN",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2006 DEC-FEB",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2007 JAN-MAR",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2007 FEB-APR",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2007 MAR-MAY",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2007 APR-JUN",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2007 MAY-JUL",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2007 JUN-AUG",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2007 JUL-SEP",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2007 AUG-OCT",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2007 SEP-NOV",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2007 OCT-DEC",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2007 NOV-JAN",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2007 DEC-FEB",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2008 JAN-MAR",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2008 FEB-APR",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2008 MAR-MAY",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2008 APR-JUN",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2008 MAY-JUL",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2008 JUN-AUG",
            "y": 72.6,
            "stringY": "72.6"
          },
          {
            "name": "2008 JUL-SEP",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "2008 AUG-OCT",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "2008 SEP-NOV",
            "y": 72.3,
            "stringY": "72.3"
          },
          {
            "name": "2008 OCT-DEC",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "2008 NOV-JAN",
            "y": 72.2,
            "stringY": "72.2"
          },
          {
            "name": "2008 DEC-FEB",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "2009 JAN-MAR",
            "y": 71.7,
            "stringY": "71.7"
          },
          {
            "name": "2009 FEB-APR",
            "y": 71.4,
            "stringY": "71.4"
          },
          {
            "name": "2009 MAR-MAY",
            "y": 71,
            "stringY": "71.0"
          },
          {
            "name": "2009 APR-JUN",
            "y": 70.8,
            "stringY": "70.8"
          },
          {
            "name": "2009 MAY-JUL",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2009 JUN-AUG",
            "y": 70.7,
            "stringY": "70.7"
          },
          {
            "name": "2009 JUL-SEP",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2009 AUG-OCT",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2009 SEP-NOV",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2009 OCT-DEC",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2009 NOV-JAN",
            "y": 70.4,
            "stringY": "70.4"
          },
          {
            "name": "2009 DEC-FEB",
            "y": 70.3,
            "stringY": "70.3"
          },
          {
            "name": "2010 JAN-MAR",
            "y": 70.2,
            "stringY": "70.2"
          },
          {
            "name": "2010 FEB-APR",
            "y": 70.2,
            "stringY": "70.2"
          },
          {
            "name": "2010 MAR-MAY",
            "y": 70.4,
            "stringY": "70.4"
          },
          {
            "name": "2010 APR-JUN",
            "y": 70.4,
            "stringY": "70.4"
          },
          {
            "name": "2010 MAY-JUL",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2010 JUN-AUG",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2010 JUL-SEP",
            "y": 70.7,
            "stringY": "70.7"
          },
          {
            "name": "2010 AUG-OCT",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2010 SEP-NOV",
            "y": 70.3,
            "stringY": "70.3"
          },
          {
            "name": "2010 OCT-DEC",
            "y": 70.4,
            "stringY": "70.4"
          },
          {
            "name": "2010 NOV-JAN",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2010 DEC-FEB",
            "y": 70.6,
            "stringY": "70.6"
          },
          {
            "name": "2011 JAN-MAR",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2011 FEB-APR",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2011 MAR-MAY",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2011 APR-JUN",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2011 MAY-JUL",
            "y": 70.2,
            "stringY": "70.2"
          },
          {
            "name": "2011 JUN-AUG",
            "y": 70.2,
            "stringY": "70.2"
          },
          {
            "name": "2011 JUL-SEP",
            "y": 70.1,
            "stringY": "70.1"
          },
          {
            "name": "2011 AUG-OCT",
            "y": 70.1,
            "stringY": "70.1"
          },
          {
            "name": "2011 SEP-NOV",
            "y": 70.1,
            "stringY": "70.1"
          },
          {
            "name": "2011 OCT-DEC",
            "y": 70.2,
            "stringY": "70.2"
          },
          {
            "name": "2011 NOV-JAN",
            "y": 70.3,
            "stringY": "70.3"
          },
          {
            "name": "2011 DEC-FEB",
            "y": 70.3,
            "stringY": "70.3"
          },
          {
            "name": "2012 JAN-MAR",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2012 FEB-APR",
            "y": 70.5,
            "stringY": "70.5"
          },
          {
            "name": "2012 MAR-MAY",
            "y": 70.7,
            "stringY": "70.7"
          },
          {
            "name": "2012 APR-JUN",
            "y": 70.9,
            "stringY": "70.9"
          },
          {
            "name": "2012 MAY-JUL",
            "y": 71.1,
            "stringY": "71.1"
          },
          {
            "name": "2012 JUN-AUG",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "2012 JUL-SEP",
            "y": 71.1,
            "stringY": "71.1"
          },
          {
            "name": "2012 AUG-OCT",
            "y": 71.1,
            "stringY": "71.1"
          },
          {
            "name": "2012 SEP-NOV",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "2012 OCT-DEC",
            "y": 71.4,
            "stringY": "71.4"
          },
          {
            "name": "2012 NOV-JAN",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "2012 DEC-FEB",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "2013 JAN-MAR",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "2013 FEB-APR",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "2013 MAR-MAY",
            "y": 71.2,
            "stringY": "71.2"
          },
          {
            "name": "2013 APR-JUN",
            "y": 71.3,
            "stringY": "71.3"
          },
          {
            "name": "2013 MAY-JUL",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "2013 JUN-AUG",
            "y": 71.5,
            "stringY": "71.5"
          },
          {
            "name": "2013 JUL-SEP",
            "y": 71.6,
            "stringY": "71.6"
          },
          {
            "name": "2013 AUG-OCT",
            "y": 71.9,
            "stringY": "71.9"
          },
          {
            "name": "2013 SEP-NOV",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "2013 OCT-DEC",
            "y": 72,
            "stringY": "72.0"
          },
          {
            "name": "2013 NOV-JAN",
            "y": 72.1,
            "stringY": "72.1"
          },
          {
            "name": "2013 DEC-FEB",
            "y": 72.4,
            "stringY": "72.4"
          },
          {
            "name": "2014 JAN-MAR",
            "y": 72.5,
            "stringY": "72.5"
          },
          {
            "name": "2014 FEB-APR",
            "y": 72.7,
            "stringY": "72.7"
          },
          {
            "name": "2014 MAR-MAY",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2014 APR-JUN",
            "y": 72.9,
            "stringY": "72.9"
          },
          {
            "name": "2014 MAY-JUL",
            "y": 72.8,
            "stringY": "72.8"
          },
          {
            "name": "2014 JUN-AUG",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2014 JUL-SEP",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2014 AUG-OCT",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2014 SEP-NOV",
            "y": 73,
            "stringY": "73.0"
          },
          {
            "name": "2014 OCT-DEC",
            "y": 73.2,
            "stringY": "73.2"
          },
          {
            "name": "2014 NOV-JAN",
            "y": 73.3,
            "stringY": "73.3"
          },
          {
            "name": "2014 DEC-FEB",
            "y": 73.4,
            "stringY": "73.4"
          },
          {
            "name": "2015 JAN-MAR",
            "y": 73.5,
            "stringY": "73.5"
          },
          {
            "name": "2015 FEB-APR",
            "y": 73.4,
            "stringY": "73.4"
          },
          {
            "name": "2015 MAR-MAY",
            "y": 73.4,
            "stringY": "73.4"
          },
          {
            "name": "2015 APR-JUN",
            "y": 73.4,
            "stringY": "73.4"
          },
          {
            "name": "2015 MAY-JUL",
            "y": 73.5,
            "stringY": "73.5"
          },
          {
            "name": "2015 JUN-AUG",
            "y": 73.6,
            "stringY": "73.6"
          },
          {
            "name": "2015 JUL-SEP",
            "y": 73.8,
            "stringY": "73.8"
          },
          {
            "name": "2015 AUG-OCT",
            "y": 73.9,
            "stringY": "73.9"
          },
          {
            "name": "2015 SEP-NOV",
            "y": 74,
            "stringY": "74.0"
          },
          {
            "name": "2015 OCT-DEC",
            "y": 74.1,
            "stringY": "74.1"
          },
          {
            "name": "2015 NOV-JAN",
            "y": 74.1,
            "stringY": "74.1"
          },
          {
            "name": "2015 DEC-FEB",
            "y": 74.1,
            "stringY": "74.1"
          },
          {
            "name": "2016 JAN-MAR",
            "y": 74.2,
            "stringY": "74.2"
          },
          {
            "name": "2016 FEB-APR",
            "y": 74.2,
            "stringY": "74.2"
          },
          {
            "name": "2016 MAR-MAY",
            "y": 74.4,
            "stringY": "74.4"
          },
          {
            "name": "2016 APR-JUN",
            "y": 74.5,
            "stringY": "74.5"
          },
          {
            "name": "2016 MAY-JUL",
            "y": 74.5,
            "stringY": "74.5"
          }
        ]
      },
      {
        "title": "Unemployment rate",
        "uri": "/employmentandlabourmarket/peoplenotinwork/unemployment/timeseries/mgsx/lms",
        "releaseDate": "Mar-May 2016",
        "latestFigure": {
          "figure": "4.9",
          "preUnit": "",
          "unit": "%"
        },
        "sparklineData": [
          {
            "name": "1971 JAN-MAR",
            "y": 3.8,
            "stringY": "3.8"
          },
          {
            "name": "1971 FEB-APR",
            "y": 3.9,
            "stringY": "3.9"
          },
          {
            "name": "1971 MAR-MAY",
            "y": 4,
            "stringY": "4.0"
          },
          {
            "name": "1971 APR-JUN",
            "y": 4.1,
            "stringY": "4.1"
          },
          {
            "name": "1971 MAY-JUL",
            "y": 4.1,
            "stringY": "4.1"
          },
          {
            "name": "1971 JUN-AUG",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "1971 JUL-SEP",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "1971 AUG-OCT",
            "y": 4.3,
            "stringY": "4.3"
          },
          {
            "name": "1971 SEP-NOV",
            "y": 4.4,
            "stringY": "4.4"
          },
          {
            "name": "1971 OCT-DEC",
            "y": 4.4,
            "stringY": "4.4"
          },
          {
            "name": "1971 NOV-JAN",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "1971 DEC-FEB",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "1972 JAN-MAR",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "1972 FEB-APR",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "1972 MAR-MAY",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "1972 APR-JUN",
            "y": 4.4,
            "stringY": "4.4"
          },
          {
            "name": "1972 MAY-JUL",
            "y": 4.4,
            "stringY": "4.4"
          },
          {
            "name": "1972 JUN-AUG",
            "y": 4.3,
            "stringY": "4.3"
          },
          {
            "name": "1972 JUL-SEP",
            "y": 4.3,
            "stringY": "4.3"
          },
          {
            "name": "1972 AUG-OCT",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "1972 SEP-NOV",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "1972 OCT-DEC",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "1972 NOV-JAN",
            "y": 4.1,
            "stringY": "4.1"
          },
          {
            "name": "1972 DEC-FEB",
            "y": 4,
            "stringY": "4.0"
          },
          {
            "name": "1973 JAN-MAR",
            "y": 3.9,
            "stringY": "3.9"
          },
          {
            "name": "1973 FEB-APR",
            "y": 3.8,
            "stringY": "3.8"
          },
          {
            "name": "1973 MAR-MAY",
            "y": 3.8,
            "stringY": "3.8"
          },
          {
            "name": "1973 APR-JUN",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "1973 MAY-JUL",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "1973 JUN-AUG",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1973 JUL-SEP",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1973 AUG-OCT",
            "y": 3.5,
            "stringY": "3.5"
          },
          {
            "name": "1973 SEP-NOV",
            "y": 3.5,
            "stringY": "3.5"
          },
          {
            "name": "1973 OCT-DEC",
            "y": 3.4,
            "stringY": "3.4"
          },
          {
            "name": "1973 NOV-JAN",
            "y": 3.4,
            "stringY": "3.4"
          },
          {
            "name": "1973 DEC-FEB",
            "y": 3.5,
            "stringY": "3.5"
          },
          {
            "name": "1974 JAN-MAR",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1974 FEB-APR",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1974 MAR-MAY",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1974 APR-JUN",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1974 MAY-JUL",
            "y": 3.6,
            "stringY": "3.6"
          },
          {
            "name": "1974 JUN-AUG",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "1974 JUL-SEP",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "1974 AUG-OCT",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "1974 SEP-NOV",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "1974 OCT-DEC",
            "y": 3.7,
            "stringY": "3.7"
          },
          {
            "name": "1974 NOV-JAN",
            "y": 3.8,
            "stringY": "3.8"
          },
          {
            "name": "1974 DEC-FEB",
            "y": 3.9,
            "stringY": "3.9"
          },
          {
            "name": "1975 JAN-MAR",
            "y": 4,
            "stringY": "4.0"
          },
          {
            "name": "1975 FEB-APR",
            "y": 4.1,
            "stringY": "4.1"
          },
          {
            "name": "1975 MAR-MAY",
            "y": 4.2,
            "stringY": "4.2"
          },
          {
            "name": "1975 APR-JUN",
            "y": 4.3,
            "stringY": "4.3"
          },
          {
            "name": "1975 MAY-JUL",
            "y": 4.5,
            "stringY": "4.5"
          },
          {
            "name": "1975 JUN-AUG",
            "y": 4.6,
            "stringY": "4.6"
          },
          {
            "name": "1975 JUL-SEP",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "1975 AUG-OCT",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "1975 SEP-NOV",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "1975 OCT-DEC",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "1975 NOV-JAN",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "1975 DEC-FEB",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "1976 JAN-MAR",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1976 FEB-APR",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1976 MAR-MAY",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1976 APR-JUN",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1976 MAY-JUL",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1976 JUN-AUG",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1976 JUL-SEP",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1976 AUG-OCT",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1976 SEP-NOV",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1976 OCT-DEC",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1976 NOV-JAN",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1976 DEC-FEB",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1977 JAN-MAR",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1977 FEB-APR",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1977 MAR-MAY",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1977 APR-JUN",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1977 MAY-JUL",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1977 JUN-AUG",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1977 JUL-SEP",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "1977 AUG-OCT",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "1977 SEP-NOV",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "1977 OCT-DEC",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "1977 NOV-JAN",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "1977 DEC-FEB",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1978 JAN-MAR",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1978 FEB-APR",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1978 MAR-MAY",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1978 APR-JUN",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1978 MAY-JUL",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1978 JUN-AUG",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1978 JUL-SEP",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1978 AUG-OCT",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1978 SEP-NOV",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1978 OCT-DEC",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1978 NOV-JAN",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1978 DEC-FEB",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1979 JAN-MAR",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1979 FEB-APR",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1979 MAR-MAY",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1979 APR-JUN",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1979 MAY-JUL",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1979 JUN-AUG",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "1979 JUL-SEP",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1979 AUG-OCT",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1979 SEP-NOV",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "1979 OCT-DEC",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "1979 NOV-JAN",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "1979 DEC-FEB",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "1980 JAN-MAR",
            "y": 5.8,
            "stringY": "5.8"
          },
          {
            "name": "1980 FEB-APR",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "1980 MAR-MAY",
            "y": 6.1,
            "stringY": "6.1"
          },
          {
            "name": "1980 APR-JUN",
            "y": 6.3,
            "stringY": "6.3"
          },
          {
            "name": "1980 MAY-JUL",
            "y": 6.5,
            "stringY": "6.5"
          },
          {
            "name": "1980 JUN-AUG",
            "y": 6.8,
            "stringY": "6.8"
          },
          {
            "name": "1980 JUL-SEP",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1980 AUG-OCT",
            "y": 7.4,
            "stringY": "7.4"
          },
          {
            "name": "1980 SEP-NOV",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "1980 OCT-DEC",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "1980 NOV-JAN",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1980 DEC-FEB",
            "y": 8.6,
            "stringY": "8.6"
          },
          {
            "name": "1981 JAN-MAR",
            "y": 8.9,
            "stringY": "8.9"
          },
          {
            "name": "1981 FEB-APR",
            "y": 9.1,
            "stringY": "9.1"
          },
          {
            "name": "1981 MAR-MAY",
            "y": 9.4,
            "stringY": "9.4"
          },
          {
            "name": "1981 APR-JUN",
            "y": 9.6,
            "stringY": "9.6"
          },
          {
            "name": "1981 MAY-JUL",
            "y": 9.7,
            "stringY": "9.7"
          },
          {
            "name": "1981 JUN-AUG",
            "y": 9.8,
            "stringY": "9.8"
          },
          {
            "name": "1981 JUL-SEP",
            "y": 9.9,
            "stringY": "9.9"
          },
          {
            "name": "1981 AUG-OCT",
            "y": 10,
            "stringY": "10.0"
          },
          {
            "name": "1981 SEP-NOV",
            "y": 10.1,
            "stringY": "10.1"
          },
          {
            "name": "1981 OCT-DEC",
            "y": 10.2,
            "stringY": "10.2"
          },
          {
            "name": "1981 NOV-JAN",
            "y": 10.3,
            "stringY": "10.3"
          },
          {
            "name": "1981 DEC-FEB",
            "y": 10.4,
            "stringY": "10.4"
          },
          {
            "name": "1982 JAN-MAR",
            "y": 10.4,
            "stringY": "10.4"
          },
          {
            "name": "1982 FEB-APR",
            "y": 10.5,
            "stringY": "10.5"
          },
          {
            "name": "1982 MAR-MAY",
            "y": 10.5,
            "stringY": "10.5"
          },
          {
            "name": "1982 APR-JUN",
            "y": 10.6,
            "stringY": "10.6"
          },
          {
            "name": "1982 MAY-JUL",
            "y": 10.6,
            "stringY": "10.6"
          },
          {
            "name": "1982 JUN-AUG",
            "y": 10.7,
            "stringY": "10.7"
          },
          {
            "name": "1982 JUL-SEP",
            "y": 10.8,
            "stringY": "10.8"
          },
          {
            "name": "1982 AUG-OCT",
            "y": 10.9,
            "stringY": "10.9"
          },
          {
            "name": "1982 SEP-NOV",
            "y": 11,
            "stringY": "11.0"
          },
          {
            "name": "1982 OCT-DEC",
            "y": 11.1,
            "stringY": "11.1"
          },
          {
            "name": "1982 NOV-JAN",
            "y": 11.2,
            "stringY": "11.2"
          },
          {
            "name": "1982 DEC-FEB",
            "y": 11.2,
            "stringY": "11.2"
          },
          {
            "name": "1983 JAN-MAR",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1983 FEB-APR",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1983 MAR-MAY",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1983 APR-JUN",
            "y": 11.4,
            "stringY": "11.4"
          },
          {
            "name": "1983 MAY-JUL",
            "y": 11.5,
            "stringY": "11.5"
          },
          {
            "name": "1983 JUN-AUG",
            "y": 11.5,
            "stringY": "11.5"
          },
          {
            "name": "1983 JUL-SEP",
            "y": 11.5,
            "stringY": "11.5"
          },
          {
            "name": "1983 AUG-OCT",
            "y": 11.6,
            "stringY": "11.6"
          },
          {
            "name": "1983 SEP-NOV",
            "y": 11.6,
            "stringY": "11.6"
          },
          {
            "name": "1983 OCT-DEC",
            "y": 11.7,
            "stringY": "11.7"
          },
          {
            "name": "1983 NOV-JAN",
            "y": 11.7,
            "stringY": "11.7"
          },
          {
            "name": "1983 DEC-FEB",
            "y": 11.8,
            "stringY": "11.8"
          },
          {
            "name": "1984 JAN-MAR",
            "y": 11.8,
            "stringY": "11.8"
          },
          {
            "name": "1984 FEB-APR",
            "y": 11.9,
            "stringY": "11.9"
          },
          {
            "name": "1984 MAR-MAY",
            "y": 11.9,
            "stringY": "11.9"
          },
          {
            "name": "1984 APR-JUN",
            "y": 11.9,
            "stringY": "11.9"
          },
          {
            "name": "1984 MAY-JUL",
            "y": 11.8,
            "stringY": "11.8"
          },
          {
            "name": "1984 JUN-AUG",
            "y": 11.8,
            "stringY": "11.8"
          },
          {
            "name": "1984 JUL-SEP",
            "y": 11.7,
            "stringY": "11.7"
          },
          {
            "name": "1984 AUG-OCT",
            "y": 11.7,
            "stringY": "11.7"
          },
          {
            "name": "1984 SEP-NOV",
            "y": 11.7,
            "stringY": "11.7"
          },
          {
            "name": "1984 OCT-DEC",
            "y": 11.6,
            "stringY": "11.6"
          },
          {
            "name": "1984 NOV-JAN",
            "y": 11.6,
            "stringY": "11.6"
          },
          {
            "name": "1984 DEC-FEB",
            "y": 11.5,
            "stringY": "11.5"
          },
          {
            "name": "1985 JAN-MAR",
            "y": 11.5,
            "stringY": "11.5"
          },
          {
            "name": "1985 FEB-APR",
            "y": 11.4,
            "stringY": "11.4"
          },
          {
            "name": "1985 MAR-MAY",
            "y": 11.4,
            "stringY": "11.4"
          },
          {
            "name": "1985 APR-JUN",
            "y": 11.4,
            "stringY": "11.4"
          },
          {
            "name": "1985 MAY-JUL",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1985 JUN-AUG",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1985 JUL-SEP",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1985 AUG-OCT",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1985 SEP-NOV",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1985 OCT-DEC",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1985 NOV-JAN",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1985 DEC-FEB",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 JAN-MAR",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 FEB-APR",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 MAR-MAY",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 APR-JUN",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 MAY-JUL",
            "y": 11.4,
            "stringY": "11.4"
          },
          {
            "name": "1986 JUN-AUG",
            "y": 11.4,
            "stringY": "11.4"
          },
          {
            "name": "1986 JUL-SEP",
            "y": 11.4,
            "stringY": "11.4"
          },
          {
            "name": "1986 AUG-OCT",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 SEP-NOV",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 OCT-DEC",
            "y": 11.3,
            "stringY": "11.3"
          },
          {
            "name": "1986 NOV-JAN",
            "y": 11.2,
            "stringY": "11.2"
          },
          {
            "name": "1986 DEC-FEB",
            "y": 11.2,
            "stringY": "11.2"
          },
          {
            "name": "1987 JAN-MAR",
            "y": 11.1,
            "stringY": "11.1"
          },
          {
            "name": "1987 FEB-APR",
            "y": 11,
            "stringY": "11.0"
          },
          {
            "name": "1987 MAR-MAY",
            "y": 10.9,
            "stringY": "10.9"
          },
          {
            "name": "1987 APR-JUN",
            "y": 10.7,
            "stringY": "10.7"
          },
          {
            "name": "1987 MAY-JUL",
            "y": 10.6,
            "stringY": "10.6"
          },
          {
            "name": "1987 JUN-AUG",
            "y": 10.4,
            "stringY": "10.4"
          },
          {
            "name": "1987 JUL-SEP",
            "y": 10.2,
            "stringY": "10.2"
          },
          {
            "name": "1987 AUG-OCT",
            "y": 10,
            "stringY": "10.0"
          },
          {
            "name": "1987 SEP-NOV",
            "y": 9.9,
            "stringY": "9.9"
          },
          {
            "name": "1987 OCT-DEC",
            "y": 9.7,
            "stringY": "9.7"
          },
          {
            "name": "1987 NOV-JAN",
            "y": 9.5,
            "stringY": "9.5"
          },
          {
            "name": "1987 DEC-FEB",
            "y": 9.4,
            "stringY": "9.4"
          },
          {
            "name": "1988 JAN-MAR",
            "y": 9.2,
            "stringY": "9.2"
          },
          {
            "name": "1988 FEB-APR",
            "y": 9,
            "stringY": "9.0"
          },
          {
            "name": "1988 MAR-MAY",
            "y": 8.9,
            "stringY": "8.9"
          },
          {
            "name": "1988 APR-JUN",
            "y": 8.7,
            "stringY": "8.7"
          },
          {
            "name": "1988 MAY-JUL",
            "y": 8.6,
            "stringY": "8.6"
          },
          {
            "name": "1988 JUN-AUG",
            "y": 8.5,
            "stringY": "8.5"
          },
          {
            "name": "1988 JUL-SEP",
            "y": 8.4,
            "stringY": "8.4"
          },
          {
            "name": "1988 AUG-OCT",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1988 SEP-NOV",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "1988 OCT-DEC",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "1988 NOV-JAN",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "1988 DEC-FEB",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "1989 JAN-MAR",
            "y": 7.6,
            "stringY": "7.6"
          },
          {
            "name": "1989 FEB-APR",
            "y": 7.4,
            "stringY": "7.4"
          },
          {
            "name": "1989 MAR-MAY",
            "y": 7.3,
            "stringY": "7.3"
          },
          {
            "name": "1989 APR-JUN",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "1989 MAY-JUL",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "1989 JUN-AUG",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1989 JUL-SEP",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1989 AUG-OCT",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1989 SEP-NOV",
            "y": 7,
            "stringY": "7.0"
          },
          {
            "name": "1989 OCT-DEC",
            "y": 7,
            "stringY": "7.0"
          },
          {
            "name": "1989 NOV-JAN",
            "y": 7,
            "stringY": "7.0"
          },
          {
            "name": "1989 DEC-FEB",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1990 JAN-MAR",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1990 FEB-APR",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1990 MAR-MAY",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1990 APR-JUN",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1990 MAY-JUL",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "1990 JUN-AUG",
            "y": 7,
            "stringY": "7.0"
          },
          {
            "name": "1990 JUL-SEP",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1990 AUG-OCT",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "1990 SEP-NOV",
            "y": 7.3,
            "stringY": "7.3"
          },
          {
            "name": "1990 OCT-DEC",
            "y": 7.5,
            "stringY": "7.5"
          },
          {
            "name": "1990 NOV-JAN",
            "y": 7.6,
            "stringY": "7.6"
          },
          {
            "name": "1990 DEC-FEB",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "1991 JAN-MAR",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "1991 FEB-APR",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "1991 MAR-MAY",
            "y": 8.5,
            "stringY": "8.5"
          },
          {
            "name": "1991 APR-JUN",
            "y": 8.7,
            "stringY": "8.7"
          },
          {
            "name": "1991 MAY-JUL",
            "y": 8.8,
            "stringY": "8.8"
          },
          {
            "name": "1991 JUN-AUG",
            "y": 9,
            "stringY": "9.0"
          },
          {
            "name": "1991 JUL-SEP",
            "y": 9.2,
            "stringY": "9.2"
          },
          {
            "name": "1991 AUG-OCT",
            "y": 9.3,
            "stringY": "9.3"
          },
          {
            "name": "1991 SEP-NOV",
            "y": 9.4,
            "stringY": "9.4"
          },
          {
            "name": "1991 OCT-DEC",
            "y": 9.5,
            "stringY": "9.5"
          },
          {
            "name": "1991 NOV-JAN",
            "y": 9.5,
            "stringY": "9.5"
          },
          {
            "name": "1991 DEC-FEB",
            "y": 9.6,
            "stringY": "9.6"
          },
          {
            "name": "1992 JAN-MAR",
            "y": 9.7,
            "stringY": "9.7"
          },
          {
            "name": "1992 FEB-APR",
            "y": 9.8,
            "stringY": "9.8"
          },
          {
            "name": "1992 MAR-MAY",
            "y": 9.8,
            "stringY": "9.8"
          },
          {
            "name": "1992 APR-JUN",
            "y": 9.8,
            "stringY": "9.8"
          },
          {
            "name": "1992 MAY-JUL",
            "y": 9.9,
            "stringY": "9.9"
          },
          {
            "name": "1992 JUN-AUG",
            "y": 9.9,
            "stringY": "9.9"
          },
          {
            "name": "1992 JUL-SEP",
            "y": 9.9,
            "stringY": "9.9"
          },
          {
            "name": "1992 AUG-OCT",
            "y": 10.1,
            "stringY": "10.1"
          },
          {
            "name": "1992 SEP-NOV",
            "y": 10.2,
            "stringY": "10.2"
          },
          {
            "name": "1992 OCT-DEC",
            "y": 10.4,
            "stringY": "10.4"
          },
          {
            "name": "1992 NOV-JAN",
            "y": 10.5,
            "stringY": "10.5"
          },
          {
            "name": "1992 DEC-FEB",
            "y": 10.7,
            "stringY": "10.7"
          },
          {
            "name": "1993 JAN-MAR",
            "y": 10.6,
            "stringY": "10.6"
          },
          {
            "name": "1993 FEB-APR",
            "y": 10.6,
            "stringY": "10.6"
          },
          {
            "name": "1993 MAR-MAY",
            "y": 10.5,
            "stringY": "10.5"
          },
          {
            "name": "1993 APR-JUN",
            "y": 10.4,
            "stringY": "10.4"
          },
          {
            "name": "1993 MAY-JUL",
            "y": 10.4,
            "stringY": "10.4"
          },
          {
            "name": "1993 JUN-AUG",
            "y": 10.3,
            "stringY": "10.3"
          },
          {
            "name": "1993 JUL-SEP",
            "y": 10.2,
            "stringY": "10.2"
          },
          {
            "name": "1993 AUG-OCT",
            "y": 10.3,
            "stringY": "10.3"
          },
          {
            "name": "1993 SEP-NOV",
            "y": 10.2,
            "stringY": "10.2"
          },
          {
            "name": "1993 OCT-DEC",
            "y": 10.3,
            "stringY": "10.3"
          },
          {
            "name": "1993 NOV-JAN",
            "y": 10.3,
            "stringY": "10.3"
          },
          {
            "name": "1993 DEC-FEB",
            "y": 10.1,
            "stringY": "10.1"
          },
          {
            "name": "1994 JAN-MAR",
            "y": 9.9,
            "stringY": "9.9"
          },
          {
            "name": "1994 FEB-APR",
            "y": 9.8,
            "stringY": "9.8"
          },
          {
            "name": "1994 MAR-MAY",
            "y": 9.7,
            "stringY": "9.7"
          },
          {
            "name": "1994 APR-JUN",
            "y": 9.7,
            "stringY": "9.7"
          },
          {
            "name": "1994 MAY-JUL",
            "y": 9.6,
            "stringY": "9.6"
          },
          {
            "name": "1994 JUN-AUG",
            "y": 9.5,
            "stringY": "9.5"
          },
          {
            "name": "1994 JUL-SEP",
            "y": 9.4,
            "stringY": "9.4"
          },
          {
            "name": "1994 AUG-OCT",
            "y": 9.3,
            "stringY": "9.3"
          },
          {
            "name": "1994 SEP-NOV",
            "y": 9.1,
            "stringY": "9.1"
          },
          {
            "name": "1994 OCT-DEC",
            "y": 9,
            "stringY": "9.0"
          },
          {
            "name": "1994 NOV-JAN",
            "y": 8.9,
            "stringY": "8.9"
          },
          {
            "name": "1994 DEC-FEB",
            "y": 8.9,
            "stringY": "8.9"
          },
          {
            "name": "1995 JAN-MAR",
            "y": 8.9,
            "stringY": "8.9"
          },
          {
            "name": "1995 FEB-APR",
            "y": 8.8,
            "stringY": "8.8"
          },
          {
            "name": "1995 MAR-MAY",
            "y": 8.8,
            "stringY": "8.8"
          },
          {
            "name": "1995 APR-JUN",
            "y": 8.7,
            "stringY": "8.7"
          },
          {
            "name": "1995 MAY-JUL",
            "y": 8.7,
            "stringY": "8.7"
          },
          {
            "name": "1995 JUN-AUG",
            "y": 8.6,
            "stringY": "8.6"
          },
          {
            "name": "1995 JUL-SEP",
            "y": 8.6,
            "stringY": "8.6"
          },
          {
            "name": "1995 AUG-OCT",
            "y": 8.6,
            "stringY": "8.6"
          },
          {
            "name": "1995 SEP-NOV",
            "y": 8.6,
            "stringY": "8.6"
          },
          {
            "name": "1995 OCT-DEC",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1995 NOV-JAN",
            "y": 8.4,
            "stringY": "8.4"
          },
          {
            "name": "1995 DEC-FEB",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1996 JAN-MAR",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "1996 FEB-APR",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1996 MAR-MAY",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1996 APR-JUN",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "1996 MAY-JUL",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "1996 JUN-AUG",
            "y": 8.1,
            "stringY": "8.1"
          },
          {
            "name": "1996 JUL-SEP",
            "y": 8.1,
            "stringY": "8.1"
          },
          {
            "name": "1996 AUG-OCT",
            "y": 8.1,
            "stringY": "8.1"
          },
          {
            "name": "1996 SEP-NOV",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "1996 OCT-DEC",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "1996 NOV-JAN",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "1996 DEC-FEB",
            "y": 7.5,
            "stringY": "7.5"
          },
          {
            "name": "1997 JAN-MAR",
            "y": 7.3,
            "stringY": "7.3"
          },
          {
            "name": "1997 FEB-APR",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "1997 MAR-MAY",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "1997 APR-JUN",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "1997 MAY-JUL",
            "y": 7.3,
            "stringY": "7.3"
          },
          {
            "name": "1997 JUN-AUG",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "1997 JUL-SEP",
            "y": 6.8,
            "stringY": "6.8"
          },
          {
            "name": "1997 AUG-OCT",
            "y": 6.7,
            "stringY": "6.7"
          },
          {
            "name": "1997 SEP-NOV",
            "y": 6.6,
            "stringY": "6.6"
          },
          {
            "name": "1997 OCT-DEC",
            "y": 6.5,
            "stringY": "6.5"
          },
          {
            "name": "1997 NOV-JAN",
            "y": 6.4,
            "stringY": "6.4"
          },
          {
            "name": "1997 DEC-FEB",
            "y": 6.4,
            "stringY": "6.4"
          },
          {
            "name": "1998 JAN-MAR",
            "y": 6.4,
            "stringY": "6.4"
          },
          {
            "name": "1998 FEB-APR",
            "y": 6.3,
            "stringY": "6.3"
          },
          {
            "name": "1998 MAR-MAY",
            "y": 6.3,
            "stringY": "6.3"
          },
          {
            "name": "1998 APR-JUN",
            "y": 6.3,
            "stringY": "6.3"
          },
          {
            "name": "1998 MAY-JUL",
            "y": 6.3,
            "stringY": "6.3"
          },
          {
            "name": "1998 JUN-AUG",
            "y": 6.3,
            "stringY": "6.3"
          },
          {
            "name": "1998 JUL-SEP",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "1998 AUG-OCT",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "1998 SEP-NOV",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "1998 OCT-DEC",
            "y": 6.1,
            "stringY": "6.1"
          },
          {
            "name": "1998 NOV-JAN",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "1998 DEC-FEB",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "1999 JAN-MAR",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "1999 FEB-APR",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "1999 MAR-MAY",
            "y": 6.1,
            "stringY": "6.1"
          },
          {
            "name": "1999 APR-JUN",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "1999 MAY-JUL",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "1999 JUN-AUG",
            "y": 5.9,
            "stringY": "5.9"
          },
          {
            "name": "1999 JUL-SEP",
            "y": 5.9,
            "stringY": "5.9"
          },
          {
            "name": "1999 AUG-OCT",
            "y": 5.8,
            "stringY": "5.8"
          },
          {
            "name": "1999 SEP-NOV",
            "y": 5.8,
            "stringY": "5.8"
          },
          {
            "name": "1999 OCT-DEC",
            "y": 5.8,
            "stringY": "5.8"
          },
          {
            "name": "1999 NOV-JAN",
            "y": 5.9,
            "stringY": "5.9"
          },
          {
            "name": "1999 DEC-FEB",
            "y": 5.8,
            "stringY": "5.8"
          },
          {
            "name": "2000 JAN-MAR",
            "y": 5.8,
            "stringY": "5.8"
          },
          {
            "name": "2000 FEB-APR",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "2000 MAR-MAY",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "2000 APR-JUN",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2000 MAY-JUL",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2000 JUN-AUG",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2000 JUL-SEP",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2000 AUG-OCT",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "2000 SEP-NOV",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2000 OCT-DEC",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2000 NOV-JAN",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2000 DEC-FEB",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2001 JAN-MAR",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2001 FEB-APR",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2001 MAR-MAY",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "2001 APR-JUN",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2001 MAY-JUL",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2001 JUN-AUG",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2001 JUL-SEP",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2001 AUG-OCT",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2001 SEP-NOV",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2001 OCT-DEC",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2001 NOV-JAN",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2001 DEC-FEB",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2002 JAN-MAR",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 FEB-APR",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 MAR-MAY",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 APR-JUN",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 MAY-JUL",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 JUN-AUG",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 JUL-SEP",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2002 AUG-OCT",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 SEP-NOV",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2002 OCT-DEC",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2002 NOV-JAN",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2002 DEC-FEB",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2003 JAN-MAR",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2003 FEB-APR",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2003 MAR-MAY",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2003 APR-JUN",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "2003 MAY-JUL",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2003 JUN-AUG",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2003 JUL-SEP",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2003 AUG-OCT",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2003 SEP-NOV",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "2003 OCT-DEC",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "2003 NOV-JAN",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2003 DEC-FEB",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2004 JAN-MAR",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2004 FEB-APR",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2004 MAR-MAY",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2004 APR-JUN",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2004 MAY-JUL",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2004 JUN-AUG",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2004 JUL-SEP",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2004 AUG-OCT",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2004 SEP-NOV",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2004 OCT-DEC",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2004 NOV-JAN",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2004 DEC-FEB",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2005 JAN-MAR",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2005 FEB-APR",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2005 MAR-MAY",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2005 APR-JUN",
            "y": 4.8,
            "stringY": "4.8"
          },
          {
            "name": "2005 MAY-JUL",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2005 JUN-AUG",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2005 JUL-SEP",
            "y": 4.7,
            "stringY": "4.7"
          },
          {
            "name": "2005 AUG-OCT",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "2005 SEP-NOV",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2005 OCT-DEC",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2005 NOV-JAN",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2005 DEC-FEB",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2006 JAN-MAR",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2006 FEB-APR",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2006 MAR-MAY",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "2006 APR-JUN",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2006 MAY-JUL",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2006 JUN-AUG",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2006 JUL-SEP",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2006 AUG-OCT",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2006 SEP-NOV",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "2006 OCT-DEC",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2006 NOV-JAN",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2006 DEC-FEB",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2007 JAN-MAR",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2007 FEB-APR",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2007 MAR-MAY",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "2007 APR-JUN",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "2007 MAY-JUL",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2007 JUN-AUG",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2007 JUL-SEP",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2007 AUG-OCT",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2007 SEP-NOV",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2007 OCT-DEC",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2007 NOV-JAN",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2007 DEC-FEB",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2008 JAN-MAR",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2008 FEB-APR",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2008 MAR-MAY",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2008 APR-JUN",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "2008 MAY-JUL",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2008 JUN-AUG",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "2008 JUL-SEP",
            "y": 5.9,
            "stringY": "5.9"
          },
          {
            "name": "2008 AUG-OCT",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "2008 SEP-NOV",
            "y": 6.2,
            "stringY": "6.2"
          },
          {
            "name": "2008 OCT-DEC",
            "y": 6.4,
            "stringY": "6.4"
          },
          {
            "name": "2008 NOV-JAN",
            "y": 6.5,
            "stringY": "6.5"
          },
          {
            "name": "2008 DEC-FEB",
            "y": 6.7,
            "stringY": "6.7"
          },
          {
            "name": "2009 JAN-MAR",
            "y": 7.1,
            "stringY": "7.1"
          },
          {
            "name": "2009 FEB-APR",
            "y": 7.3,
            "stringY": "7.3"
          },
          {
            "name": "2009 MAR-MAY",
            "y": 7.6,
            "stringY": "7.6"
          },
          {
            "name": "2009 APR-JUN",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2009 MAY-JUL",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2009 JUN-AUG",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2009 JUL-SEP",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2009 AUG-OCT",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2009 SEP-NOV",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2009 OCT-DEC",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2009 NOV-JAN",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "2009 DEC-FEB",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2010 JAN-MAR",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "2010 FEB-APR",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "2010 MAR-MAY",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2010 APR-JUN",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2010 MAY-JUL",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2010 JUN-AUG",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2010 JUL-SEP",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2010 AUG-OCT",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2010 SEP-NOV",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2010 OCT-DEC",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2010 NOV-JAN",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2010 DEC-FEB",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2011 JAN-MAR",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2011 FEB-APR",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "2011 MAR-MAY",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2011 APR-JUN",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2011 MAY-JUL",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "2011 JUN-AUG",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "2011 JUL-SEP",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "2011 AUG-OCT",
            "y": 8.4,
            "stringY": "8.4"
          },
          {
            "name": "2011 SEP-NOV",
            "y": 8.5,
            "stringY": "8.5"
          },
          {
            "name": "2011 OCT-DEC",
            "y": 8.4,
            "stringY": "8.4"
          },
          {
            "name": "2011 NOV-JAN",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "2011 DEC-FEB",
            "y": 8.3,
            "stringY": "8.3"
          },
          {
            "name": "2012 JAN-MAR",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "2012 FEB-APR",
            "y": 8.2,
            "stringY": "8.2"
          },
          {
            "name": "2012 MAR-MAY",
            "y": 8.1,
            "stringY": "8.1"
          },
          {
            "name": "2012 APR-JUN",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "2012 MAY-JUL",
            "y": 8.1,
            "stringY": "8.1"
          },
          {
            "name": "2012 JUN-AUG",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2012 JUL-SEP",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2012 AUG-OCT",
            "y": 7.9,
            "stringY": "7.9"
          },
          {
            "name": "2012 SEP-NOV",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2012 OCT-DEC",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2012 NOV-JAN",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2012 DEC-FEB",
            "y": 8,
            "stringY": "8.0"
          },
          {
            "name": "2013 JAN-MAR",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2013 FEB-APR",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2013 MAR-MAY",
            "y": 7.8,
            "stringY": "7.8"
          },
          {
            "name": "2013 APR-JUN",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "2013 MAY-JUL",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "2013 JUN-AUG",
            "y": 7.7,
            "stringY": "7.7"
          },
          {
            "name": "2013 JUL-SEP",
            "y": 7.6,
            "stringY": "7.6"
          },
          {
            "name": "2013 AUG-OCT",
            "y": 7.4,
            "stringY": "7.4"
          },
          {
            "name": "2013 SEP-NOV",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "2013 OCT-DEC",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "2013 NOV-JAN",
            "y": 7.2,
            "stringY": "7.2"
          },
          {
            "name": "2013 DEC-FEB",
            "y": 6.9,
            "stringY": "6.9"
          },
          {
            "name": "2014 JAN-MAR",
            "y": 6.8,
            "stringY": "6.8"
          },
          {
            "name": "2014 FEB-APR",
            "y": 6.6,
            "stringY": "6.6"
          },
          {
            "name": "2014 MAR-MAY",
            "y": 6.4,
            "stringY": "6.4"
          },
          {
            "name": "2014 APR-JUN",
            "y": 6.3,
            "stringY": "6.3"
          },
          {
            "name": "2014 MAY-JUL",
            "y": 6.1,
            "stringY": "6.1"
          },
          {
            "name": "2014 JUN-AUG",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "2014 JUL-SEP",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "2014 AUG-OCT",
            "y": 6,
            "stringY": "6.0"
          },
          {
            "name": "2014 SEP-NOV",
            "y": 5.9,
            "stringY": "5.9"
          },
          {
            "name": "2014 OCT-DEC",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "2014 NOV-JAN",
            "y": 5.7,
            "stringY": "5.7"
          },
          {
            "name": "2014 DEC-FEB",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "2015 JAN-MAR",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "2015 FEB-APR",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2015 MAR-MAY",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "2015 APR-JUN",
            "y": 5.6,
            "stringY": "5.6"
          },
          {
            "name": "2015 MAY-JUL",
            "y": 5.5,
            "stringY": "5.5"
          },
          {
            "name": "2015 JUN-AUG",
            "y": 5.4,
            "stringY": "5.4"
          },
          {
            "name": "2015 JUL-SEP",
            "y": 5.3,
            "stringY": "5.3"
          },
          {
            "name": "2015 AUG-OCT",
            "y": 5.2,
            "stringY": "5.2"
          },
          {
            "name": "2015 SEP-NOV",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2015 OCT-DEC",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2015 NOV-JAN",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2015 DEC-FEB",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2016 JAN-MAR",
            "y": 5.1,
            "stringY": "5.1"
          },
          {
            "name": "2016 FEB-APR",
            "y": 5,
            "stringY": "5.0"
          },
          {
            "name": "2016 MAR-MAY",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "2016 APR-JUN",
            "y": 4.9,
            "stringY": "4.9"
          },
          {
            "name": "2016 MAY-JUL",
            "y": 4.9,
            "stringY": "4.9"
          }
        ]
      },
      {
        "title": "United Kingdom population mid-year estimate",
        "uri": "/peoplepopulationandcommunity/populationandmigration/populationestimates/timeseries/ukpop/pop",
        "releaseDate": "Mar-May 2016",
        "latestFigure": {
          "figure": "65,110,000",
          "preUnit": "",
          "unit": ""
        },
        "sparklineData": []
      }
    ],
    "featured": [
      {
        "title": "UK perspectives 2016",
        "description": "Analysis and data to help you understand key issues in the EU referendum debate",
        "uri": "http://visual.ons.gov.uk/category/uncategorized/uk-perspectives-2016/"
      },
      {
        "title": "Data science campus",
        "description": "World-leading expertise in data science, and to benefit from faster, richer and more precise economic data",
        "uri": "/aboutus/whatwedo/datasciencecampus"
      },
      {
        "title": "Independent review of UK economic statistics",
        "description": "The final report of the independent review of UK economic statistics, led by Professor Sir Charles Bean",
        "uri": "https://www.gov.uk/government/publications/independent-review-of-uk-economic-statistics-final-report"
      }
    ],
    "other": [
      {
        "title": "Visual.ONS",
        "uri": "https://visual.ons.gov.uk",
        "description": "Making ONS statistics accessible and relevant to a wider public audience"
      },
      {
        "title": "Census",
        "uri": "/census",
        "description": "Discover how our census statistics paint a picture of the nation we live in"
      },
      {
        "title": "Looking for local statistics?",
        "uri": "/help/localstatistics",
        "description": "A handy guide on where to find local statistics"
      },
      {
        "title": "GOV.UK",
        "uri": "https://www.gov.uk/government/statistics/announcements",
        "description": "Other official statistics produced impartially and free from political influence"
      }
    ]
  }
}`

func Handler(rendererURL string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		b := []byte(stubbedData)
		rdr := bytes.NewReader(b)

		rendererReq, err := http.NewRequest("POST", rendererURL+"/homepage", rdr)
		if err != nil {
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FIXME there's other headers we want
		rendererReq.Header.Set("Accept-Language", string(lang.Get(req)))
		rendererReq.Header.Set("X-Request-Id", req.Header.Get("X-Request-Id"))

		res, err := http.DefaultClient.Do(rendererReq)
		if err != nil {
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			err = fmt.Errorf("unexpected status code: %d", res.StatusCode)
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FIXME should stream this using a io.Reader etc
		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		for hdr, v := range res.Header {
			for _, v2 := range v {
				w.Header().Add(hdr, v2)
			}
		}
		w.WriteHeader(res.StatusCode)
		w.Write(b)
	}
}
