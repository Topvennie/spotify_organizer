import { AuthLayout } from "@/layout/AuthLayout"
import { useAuth } from "@/lib/hooks/useAuth"
import { Button } from "@mantine/core"

export const Index = () => {
  const { logout } = useAuth()

  return (
    <AuthLayout>
      <div className="flex flex-col items-center pt-20 px-4">
        <h1 className="text-4xl font-bold mb-6 text-center">Hi!</h1>

        <div className="max-w-2xl text-left space-y-4">
          <p className="text-lg">Good luck with the new project!</p>
        </div>

        <Button onClick={logout}>Logout</Button>
      </div>
    </AuthLayout>
  )
}

