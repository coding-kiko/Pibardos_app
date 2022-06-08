import { Button, Flex, Text } from 'native-base'
import theme from '../../theme'
import { Image } from 'react-native'

const AuthScreen = () => {
  return (

    <Flex flex='1' >

      <Flex direction='row' justify='center' align='center' pt='50%'>
        <Image
          source={require('../assets/logos/favicon.png') }
          width='20%'
          height= '20%'
          alt='logo'>

        </Image>

        <Text pl='2' fontSize='40' fontWeight='bold'>Pibardos App</Text>
      </Flex>
      <Flex flex='1' justify='center' align='center' mt='-30%'>

      <Button width='90%' bg={theme.colors.blue} height='8%' onPress={() => console.log('login press')}>
        Iniciar Sesion
      </Button>
      <Button mt='3%' bg={theme.colors.green} width='90%' height='8%' onPress={() => console.log('register pressed')}>
        Crear una nueva cuenta
      </Button>
      </Flex>
      <Flex align='center' pb='5%'>
      <Text>Hecho con 🖤 por Frank, Mazen e Ivanchu </Text>
      </Flex>
    </Flex>

  )
}

export default AuthScreen
