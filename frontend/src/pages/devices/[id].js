import { useRouter } from 'next/router'

const Page = () => {
    const router = useRouter()
    const { id } = router.query

    return <p>Device: {id}</p>
}
    

export default Page;
