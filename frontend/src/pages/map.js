import Head from 'next/head';
import { GoogleMap, useLoadScript, Marker, OverlayView } from "@react-google-maps/api"
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useEffect, useMemo, useState } from 'react';
import mapStyles from '../utils/mapStyles.json';
import { useRouter } from 'next/router';

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

    const router = useRouter();

    const [mapRef, setMapRef] = useState(null);

    const [mapCenter, setMapCenter] = useState(null);
    const [markers, setMarkers] = useState([]);
    const [activeMarker, setActiveMarker] = useState(null);
    const [activeMarkerdata, setActiveMarkerdata] = useState(null);
    const [zoom, setZoom] = useState(null)

    const handleDragEnd = () => {
      if (mapRef) {
        const newCenter = mapRef.getCenter();
        console.log("newCenter:",newCenter.lat(), newCenter.lng());
        localStorage.setItem("mapCenter", JSON.stringify({"lat":newCenter.lat(),"lng":newCenter.lng()}))
      }
    }

    const handleZoomChange = () => {
      if (mapRef) {
        const newZoom = mapRef.getZoom();
        console.log("new zoom", newZoom)
        localStorage.setItem("zoom", newZoom)
      }
    }

    const handleOnLoad = map => {
      setMapRef(map);
    };

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
        setMarkers([])
        console.log("error to get map markers")
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

      let zoomFromLocalStorage = localStorage.getItem("zoom")
      if (zoomFromLocalStorage) {
        setZoom(Number(zoomFromLocalStorage))
      }else{
        setZoom(25)
      }

      let mapCenterFromLocalStorage = localStorage.getItem("mapCenter")
      if (mapCenterFromLocalStorage){
        let fmtMapCenter = JSON.parse(localStorage.getItem("mapCenter"))
        console.log("mapCenterFromLocalStorage:", fmtMapCenter)
        setMapCenter({
          lat: Number(fmtMapCenter.lat),
          lng: Number(fmtMapCenter.lng),
        })
        return
      }
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
            console.log("Error getting user location:", error);
          }
        );
      } else {
        // Geolocation is not supported by the browser
        console.log("Geolocation is not supported by this browser, or the user denied access");
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
    
    return ( markers && zoom &&
    <>
      <Head>
        <title>
          Maps | Oktopus
        </title>
      </Head>
      <GoogleMap
        options={mapOptions}
        zoom={zoom}
        center={mapCenter ? mapCenter : {
          lat: 0.0,
          lng: 0.0,
        }}
        mapContainerStyle={{ width: '100%', height: '100%' }}
        onLoad={handleOnLoad}
        clickableIcons={false}
        onDragEnd={handleDragEnd}
        onZoomChanged={handleZoomChange}
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