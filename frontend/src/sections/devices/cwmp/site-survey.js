import PropTypes from 'prop-types';
import ArrowPathIcon from '@heroicons/react/24/solid/ArrowPathIcon';
import ArrowRightIcon from '@heroicons/react/24/solid/ArrowRightIcon';
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
    Paper,
    Container,
    CircularProgress
} from '@mui/material';
import { Scrollbar } from 'src/components/scrollbar';
import { alpha, useTheme } from '@mui/material/styles';
import { Chart } from 'src/components/chart';
import ChartBarSquareIcon from '@heroicons/react/24/outline/ChartBarSquareIcon';
import ListBulletIcon from '@heroicons/react/24/outline/ListBulletIcon';
import { useRouter } from 'next/router';
import { Stack } from '@mui/system';
import { useEffect, useState } from 'react';

const useChartOptions = () => {
    const theme = useTheme();

    return {
        chart: {
            background: 'transparent',
            stacked: false,
            toolbar: {
                show: true
            }
        },
        colors: [
            theme.palette.graphics.dark,
            theme.palette.graphics.darkest,
            theme.palette.graphics.light,
            theme.palette.graphics.main,
            theme.palette.graphics.lightest,
        ],
        dataLabels: {
            enabled: false
        },
        fill: {
            opacity: 1,
            type: 'solid'
        },
        grid: {
            borderColor: theme.palette.divider,
            strokeDashArray: 2,
            xaxis: {
                lines: {
                    show: false
                }
            },
            yaxis: {
                lines: {
                    show: true
                }
            }
        },
        legend: {
            show: true
        },
        plotOptions: {
            bar: {
                columnWidth: '40px'
            }
        },
        stroke: {
            colors: ['transparent'],
            show: true,
            width: 2
        },
        theme: {
            mode: theme.palette.mode
        },
        xaxis: {
            axisBorder: {
                color: theme.palette.divider,
                show: true
            },
            axisTicks: {
                color: theme.palette.divider,
                show: true
            },
            categories: [
                'Jan',
                'Feb',
                'Mar',
                'Apr',
                'May',
                'Jun',
                'Jul',
                'Aug',
                'Sep',
                'Oct',
                'Nov',
                'Dec'
            ],
            labels: {
                offsetY: 5,
                style: {
                    colors: theme.palette.text.secondary
                }
            }
        },
        yaxis: {
            labels: {
                formatter: (value) => (value > 0 ? `${value}K` : `${value}`),
                offsetX: -10,
                style: {
                    colors: theme.palette.text.secondary
                }
            }
        }
    };
};

export const SiteSurvey = (props) => {
    const chartSeries = [{ name: 'This year', data: [18, 16, 5, 8, 3, 14, 14, 16, 17, 19, 18, 20] }, { name: 'Last year', data: [12, 11, 4, 6, 2, 9, 9, 10, 11, 12, 13, 13] }]

    const router = useRouter();
    const [content, setContent] = useState(null);
    const [frequency, setFrequency] = useState("2.4GHz");

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
                    <Button
                        color="inherit"
                        size="small"
                        disabled="true"
                    >
                        <SvgIcon>
                            <ChartBarSquareIcon />
                        </SvgIcon>
                    </Button>
                    <Button
                        color="inherit"
                        size="small"
                    >
                        <SvgIcon>
                            <ListBulletIcon />
                        </SvgIcon>
                    </Button>
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
            <Card sx={{ height: '100%' }}>
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
                {/* <CardContent>
                <Chart
                height={500}
                options={chartOptions}
                series={chartSeries}
                type="bar"
                width="100%"
                />
            </CardContent>
            <Divider />
            <CardActions sx={{ justifyContent: 'center' }}>
            <FormControl xs={2}>
                <RadioGroup
                    aria-labelledby="demo-controlled-radio-buttons-group"
                    name="controlled-radio-buttons-group"
                    value={"2.4GHz"}
                >
                <Grid container>
                    <FormControlLabel value="2.4GHz" control={<Radio />} label="2.4GHz" />
                    <FormControlLabel value="5GHz" control={<Radio />} label="5GHz" />
                </Grid>
            </RadioGroup>
            </FormControl>
            </CardActions> */}
            </Card>
        </Grid>: <CircularProgress></CircularProgress>}
        </Stack>
    );
};