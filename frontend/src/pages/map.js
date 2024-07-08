import Head from 'next/head';
import { GoogleMap, useLoadScript, Marker, OverlayView } from "@react-google-maps/api"
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useEffect, useMemo, useState } from 'react';
import mapStyles from '../utils/mapStyles.json';

const getPixelPositionOffset = pixelOffset => (width, height) => ({
  x: -(width / 2) + pixelOffset.x,
  y: -(height / 2) + pixelOffset.y
});

const Popup = props => {
  return (
    <OverlayView
      position={props.anchorPosition}
      mapPaneName={OverlayView.OVERLAY_MOUSE_TARGET}
      getPixelPositionOffset={getPixelPositionOffset(props.markerPixelOffset)}
    >
      <div className="popup-tip-anchor">
        <div className="popup-bubble-anchor">
          <div className="popup-bubble-content">{props.content}</div>
        </div>
      </div>
    </OverlayView>
  );
};

const Page = () => {

    const libraries = useMemo(() => ['places'], []);

    const [mapCenter, setMapCenter] = useState(null);
    const [markers, setMarkers] = useState([]);
    const [activeMarker, setActiveMarker] = useState(null);
    const [activeMarkerdata, setActiveMarkerdata] = useState(null);

    const fetchMarkers = async () => {

      var myHeaders = new Headers();
      myHeaders.append("Content-Type", "application/json");
      myHeaders.append("Authorization", localStorage.getItem("token"));
  
      var requestOptions = {
        method: 'GET',
        headers: myHeaders,
        redirect: 'follow'
      };

      let result = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/map`, requestOptions)

      if (result.status == 200) {
        const content = await result.json()
        setMarkers(content)
      }else if (result.status == 403) {
        console.log("num tenx permissão, seu boca de sandália")
        return router.push("/403")
      }else if (result.status == 401){
        console.log("taix nem autenticado, sai fora oh")
        return router.push("/auth/login")
      } else {
        console.log("agora quebrasse ux córno mô quiridu")
        const content = await result.json()
        throw new Error(content);
      }    
    }

    const fetchActiveMarkerData = async (id) => {
      var myHeaders = new Headers();
      myHeaders.append("Content-Type", "application/json");
      myHeaders.append("Authorization", localStorage.getItem("token"));
  
      var requestOptions = {
        method: 'GET',
        headers: myHeaders,
        redirect: 'follow'
      };

      let result = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?id=`+id, requestOptions)

      if (result.status == 200) {
        const content = await result.json()
        setActiveMarkerdata(content)
      }else if (result.status == 403) {
        return router.push("/403")
      }else if (result.status == 401){
        return router.push("/auth/login")
      } else {
        console.log("no device info found")
        const content = await result.json()
      }    
    }

    useEffect(()=> {
      fetchMarkers();
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
    
    return ( mapCenter && markers &&
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
      >
        {
          markers.map((marker, index) => (
            <Marker
              key={index}
              position={{ lat: marker.coordinates.lat, lng: marker.coordinates.lng }}
              icon={{
                url: marker.img,
                scaledSize: new window.google.maps.Size(50, 50),
                anchor: new window.google.maps.Point(25, 25),
              }}
              draggable={false}
              clickable={true}
              onClick={() => {
                setActiveMarkerdata(null);
                if (activeMarker?.sn === marker.sn) {
                  setActiveMarker(null);
                  return;
                }
                fetchActiveMarkerData(marker.sn);
                setActiveMarker({
                  sn: marker.sn,
                  position: { lat: marker.coordinates.lat, lng: marker.coordinates.lng }
                });
              }}
            >
            </Marker>
          ))
        }
        {activeMarker &&
        <Popup
          anchorPosition={activeMarker.position}
          markerPixelOffset={{ x: 0, y: -32 }}
          content={activeMarkerdata  ?
                <div>
                  <div>SN: {activeMarker.sn}</div>
                  <div>
                    <div>Model: {activeMarkerdata.Model?activeMarkerdata.Model:activeMarkerdata.ProductClass}</div>
                    <div>Alias: {activeMarkerdata.Alias}</div>
                    <div>Status: {activeMarkerdata.Status == 2 ? <span style={{color:"green"}}>online</span> : <span style={{color:"red"}}>offline</span>}</div>
                  </div>
                </div>
              : <p>no device info found</p>}
        />}
      </GoogleMap>
    </>
  )};
  
  Page.getLayout = (page) => (
    <DashboardLayout>
      {page}
    </DashboardLayout>
  );
  
  export default Page;