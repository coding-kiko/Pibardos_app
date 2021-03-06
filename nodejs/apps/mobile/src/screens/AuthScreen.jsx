import { Box, Text, Link } from 'native-base'

const AuthScreen = () => {
  return (
    <Box>
      <Text mx='16'>
        NativeBase is a component library that enables devs to build universal
        design systems. It is built on top of React Native, allowing you to
        develop apps for Android, iOS and the Web.{' '}
        <Link href='https://nativebase.io' isExternal _text={{
          color: 'blue.400'
        }} mt={-0.5} _web={{
          mb: -2
        }}>
          Read More
        </Link>
      </Text>
    </Box>
  )
}

export default AuthScreen
