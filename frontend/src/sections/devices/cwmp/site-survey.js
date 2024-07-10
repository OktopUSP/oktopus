import {
    Button,
    Card,
    CardActions,
    CardContent,
    CardHeader,
    Divider,
    SvgIcon,
    FormControl,
    FormLabel,
    Radio,
    RadioGroup,
    Grid,
    FormControlLabel,
    Box,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow,
    TableContainer,
    CircularProgress,
    ToggleButton,
    ToggleButtonGroup
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import { Chart } from 'src/components/chart';
import ChartBarSquareIcon from '@heroicons/react/24/outline/ChartBarSquareIcon';
import ListBulletIcon from '@heroicons/react/24/outline/ListBulletIcon';
import { useRouter } from 'next/router';
import { Stack } from '@mui/system';
import { useEffect, useState } from 'react';

export const SiteSurvey = (props) => {
    
    // const getMaxChannel = () => {
    //     if (frequency == "2.4GHz") {
    //         return 13;
    //     } else {
    //         return 128;
    //     }
    // }

    // const geMinChannel = () => {
    //     if (frequency == "2.4GHz") {
    //         return 0;
    //     } else {
    //         return 36;
    //     }
    // }

    // const getChannelSpacing = () => {
    //     if (frequency == "2.4GHz") {
    //         return 1;
    //     } else {
    //         return 4;
    //     }
    // }

    // const getChannelAmount = () => {
    //     if (frequency == "2.4GHz") {
    //         return 13;
    //     } else {
    //         return 20;
    //     }
    // }

    const getCategories = () => {

    }

    const router = useRouter();

    const [frequency, setFrequency] = useState("2.4GHz");
    const [view, setView] = useState("chart");
    const [content, setContent] = useState(null);

    const getSeries = () => {
        let series = []
        content[frequency].map((network) => {

            let data = []

            if (frequency == "2.4GHz") {
                if (Number(network.bandwidth) == 20) {
                    data.push({"x": Number(network.channel) -2, "y": -100})
                    data.push({"x": Number(network.channel), "y": Number(network.signal_level)})
                    data.push({"x": Number(network.channel) +2, "y": -100})
                }
                if (Number(network.bandwidth) == 40) {
                    data.push({"x": Number(network.channel) -4, "y": -100})
                    data.push({"x": Number(network.channel), "y": Number(network.signal_level)})
                    data.push({"x": Number(network.channel) +4, "y": -100})
                }
            }else {
                if (Number(network.bandwidth) == 20) {
                    data.push({"x": Number(network.channel) -4, "y": -100})
                    data.push({"x": Number(network.channel), "y": Number(network.signal_level)})
                    data.push({"x": Number(network.channel) +4, "y": -100})
                }
                if (Number(network.bandwidth) == 40) {
                    data.push({"x": Number(network.channel) -8, "y": -100})
                    data.push({"x": Number(network.channel), "y": Number(network.signal_level)})
                    data.push({"x": Number(network.channel) +8, "y": -100})
                }
                if (Number(network.bandwidth) == 80) {
                    data.push({"x": Number(network.channel) -16, "y": -100})
                    data.push({"x": Number(network.channel), "y": Number(network.signal_level)})
                    data.push({"x": Number(network.channel) +16, "y": -100})
                }
                if (Number(network.bandwidth) == 160) {
                    data.push({"x": Number(network.channel) -32, "y": -100})
                    data.push({"x": Number(network.channel), "y": Number(network.signal_level)})
                    data.push({"x": Number(network.channel) +32, "y": -100})
                }
            }

            let ssid = network.ssid
            if ( ssid == "") {
                ssid = " "
            }
            return series.push({
                name: ssid,
                data: data
            })
        })
        return series;
    }

    const useChartOptions = () => {
        const theme = useTheme();
    
        return {
            chart: {
                background: 'transparent',
                stacked: false,
                toolbar: {
                    show: true
                },
                zoom: {
                    enabled: false
                },
            },
            title: {
                text: 'Site Survey Results',
            },
            // markers: {
            //     size: 5,
            //     hover: {
            //       size: 9
            //     }
            // },
            colors: [
                theme.palette.graphics.dark,
                theme.palette.warning.main,
                theme.palette.graphics.darkest,
                theme.palette.graphics.main,
                theme.palette.info.light,
                theme.palette.graphics.lightest,
                theme.palette.primary.main,
                theme.palette.graphics.light,
                theme.palette.error.light,
                theme.palette.error.dark
            ],
            dataLabels: {
                enabled: false
            },
            grid: {
                //borderColor: theme.palette.divider,
                strokeDashArray: 2,
                xaxis: {
                    lines: {
                        show: true
                    }
                },
                yaxis: {
                    lines: {
                        show: true
                    },
                },
            },
            legend: {
                show: true,
                showForSingleSeries: true,
            },
            plotOptions: {
                area: {
                    fillTo: 'end',
                }
            },
            stroke: {
                show: true,
                curve: 'smooth',
                lineCap: 'round',
            },
            theme: {
                mode: theme.palette.mode
            },
            yaxis: {
                min: -100,
                max:  0,
                labels: {
                  formatter: function (value) {
                    return value + ' dBm';
                  },
                  //offsetY: -10,
                  style: {
                      //colors: theme.palette.text.secondary
                  }
                },
            },
            // annotations: {
            //     xaxis: [
            //         {
            //         x: 9,
            //         x2: 10,
            //         borderColor: '#0b54ff',
            //         label: {
            //             style: {
            //             color: 'black',
            //             },
            //             text: 'Channel 10',
            //             offsetX: 45,
            //             borderColor: 'transparent',
            //             style: {
            //                 background: 'transparent',
            //                 color:  theme.palette.text.secondary,
            //                 fontSize: '17px',
            //             },
            //         }
            //         }
            //     ]
            // },
            // annotations: {
            //     points: [{
            //       x: 9,
            //       y: -5,
            //       label: {
            //         borderColor: '#775DD0',
            //         offsetY: 0,
            //         style: {
            //           color: '#fff',
            //           background: '#775DD0',
            //         },
            //         rotate: -45,
            //         text: 'Bananas are good',
            //       }
            //     }]
            //   },
            xaxis: {
                // tickPlacement: 'on',
                // tickAmount: getChannelAmount(),
                tickPlacement: 'on',
                labels: {
                    show: true,
                    style: {
                        //colors: theme.palette.text.secondary
                    },
                    trim: true,
                },
                // max: getMaxChannel(),
                // min: geMinChannel(),
                // stepSize: getChannelSpacing(),
                //type: 'category',
                //categories: [getCategories()],
                type: 'numeric',
                decimalsInFloat: 0,
            },
            tooltip: {
                x: {
                    show: true,
                    formatter: (seriesName) => "Channel "+ seriesName,
                },
                followCursor: false,
                intersect: false,
                shared: true,
                enabled: true,
                onDatasetHover: {
                    highlightDataSeries: true,
                }
            }
        };
    };

    const chartOptions = useChartOptions();

    const fetchSiteSurveyData = async () => {

        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", localStorage.getItem("token"));
      
        var requestOptions = {
          method: 'GET',
          headers: myHeaders,
          redirect: 'follow'
        };

        fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device/${router.query.id[0]}/sitesurvey`, requestOptions)
        .then(response => {
            if (response.status === 401) {
                router.push("/auth/login")
            }
            return response.json()
        })
        .then(result => {
            console.log("sitesurvey content", result)
            setContent(result)
        })
        .catch(error => console.log('error', error));
    };


    useEffect(() => {
        fetchSiteSurveyData();
    },[])

    return (
        <Stack 
        justifyContent="center" 
        alignItems={!content &&"center"}
        >
        {content ?
        <Grid spacing={1}>
            <Grid container>
                <Card>
                    <ToggleButtonGroup
                    value={view}
                    exclusive
                    onChange={(e, value) => {
                        setView(value)
                    }}
                    >
                        <ToggleButton
                            size="small"
                            value="chart"
                        >
                            <SvgIcon>
                                <ChartBarSquareIcon />
                            </SvgIcon>
                        </ToggleButton>
                        <ToggleButton
                            size="small"
                            value="list"
                        >
                            <SvgIcon>
                                <ListBulletIcon />
                            </SvgIcon>
                        </ToggleButton>
                    </ToggleButtonGroup>
                </Card>
            </Grid>
            <Box display="flex"
                justifyContent="center"
                alignItems="center"
                marginBottom={3}
                >
                <FormControl>
                    <RadioGroup
                        aria-labelledby="demo-controlled-radio-buttons-group"
                        name="controlled-radio-buttons-group"
                        value={frequency}
                        onChange={(e) => {
                            setFrequency(e.target.value)
                        }}
                    >
                        <Grid container>
                            <FormControlLabel value="2.4GHz" control={<Radio />} label="2.4GHz" />
                            <FormControlLabel value="5GHz" control={<Radio />} label="5GHz" />
                        </Grid>
                    </RadioGroup>
                </FormControl>
            </Box>
            {view == "list" && <Card sx={{ height: '100%' }}>
                <Box sx={{ minWidth: 800, }}>
                    <TableContainer sx={{ maxHeight: 600 }}>
                        <Table exportButton={true}>
                            <TableHead>
                                <TableRow>
                                    <TableCell align="center">
                                        SSID
                                    </TableCell>
                                    <TableCell>
                                        Channel
                                    </TableCell>
                                    <TableCell>
                                        BandWidth
                                    </TableCell>
                                    <TableCell>
                                        Standard
                                    </TableCell>
                                    <TableCell>
                                        Signal Level
                                    </TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {content[frequency] && content[frequency].map((c) => {
                                    return (
                                        <TableRow
                                            hover
                                            key={c.ssid + c.signal_level}
                                        >
                                            <TableCell align="center">
                                                {c.ssid}
                                            </TableCell>
                                            <TableCell>
                                                {c.channel}
                                            </TableCell>
                                            <TableCell>
                                                {c.bandwidth}MHz
                                            </TableCell>
                                            <TableCell>
                                                {c.standard}
                                            </TableCell>
                                            <TableCell>
                                                {c.signal_level} dbm
                                            </TableCell>
                                        </TableRow>
                                    );
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Box>
            </Card>}
            {view == "chart" && <Card>
                <CardContent>
                    <Chart
                    height={500}
                    options={chartOptions}
                    series={getSeries()}
                    type="area"
                    width="100%"
                    />
                </CardContent>
            </Card>}
        </Grid>: <CircularProgress></CircularProgress>}
        </Stack>
    );
};