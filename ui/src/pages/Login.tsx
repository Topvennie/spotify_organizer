import { useAuth } from "@/lib/hooks/useAuth";
import { Button, Center, Paper, Stack, Text, Title } from "@mantine/core";
import { FaSpotify } from "react-icons/fa6";

export const Login = () => {
  const { login } = useAuth();

  return (
    <Center h="100vh">
      <Paper shadow="sm" p="xl" className="w-96">
        <Stack align="center">
          <Title>Login</Title>
          <Text c="gray">Spotify Organizer</Text>
          <Button onClick={login} color="red" size="lg" className="my-12">
            <FaSpotify size={"1.7rem"} className="mr-2" />
            Spotify
          </Button>
        </Stack>
      </Paper>
    </Center>
  );
};
