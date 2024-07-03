import Head from 'next/head';
import { GoogleMap, useLoadScript } from "@react-google-maps/api"
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useEffect, useMemo, useState } from 'react';
import mapStyles from '../utils/mapStyles.json';

const Page = () => {

    const libraries = useMemo(() => ['places'], []);

    const [mapCenter, setMapCenter] = useState(null);

    useEffect(()=> {
        // Check if geolocation is supported by the browser
        if ("geolocation" in navigator) {
          // Prompt user for permission to access their location
          navigator.geolocation.getCurrentPosition(
            // Get the user's latitude and longitude coordinates
            // Success callback function
            function(position) {
              // Update the map with the user's new location
              setMapCenter({
                lat: position.coords.latitude,
                lng: position.coords.longitude,
              })
            },
            // Error callback function
            function(error) {
              // Handle errors, e.g. user denied location sharing permissions
              console.error("Error getting user location:", error);
            }
          );
        } else {
          // Geolocation is not supported by the browser
          console.error("Geolocation is not supported by this browser.");
        }
    },[])
  
    const mapOptions = useMemo(
      () => ({
        disableDefaultUI: false,
        clickableIcons: true,
        zoomControl: true,
        controlSize: 23,
        styles: mapStyles,
        mapTypeControlOptions: {
          mapTypeIds: ['roadmap', 'satellite'],
        }
      }),
      []
    );
  
    const { isLoaded } = useLoadScript({
      googleMapsApiKey: process.env.NEXT_PUBLIC_GOOGLE_MAPS_KEY,
      libraries: libraries,
    });
  
    if (!isLoaded) {
      return <p>Loading...</p>;
    }
    
    return ( mapCenter &&
    <>
      <Head>
        <title>
          Maps | Oktopus
        </title>
      </Head>
      <GoogleMap
        options={mapOptions}
        zoom={14}
        center={mapCenter}
        mapContainerStyle={{ width: '100%', height: '100%' }}
        onLoad={() => console.log('Map Component Loaded...')}
        clickableIcons={false}
      />
    </>
  )};
  
  Page.getLayout = (page) => (
    <DashboardLayout>
      {page}
    </DashboardLayout>
  );
  
  export default Page;