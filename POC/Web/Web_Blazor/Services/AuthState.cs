using System;

namespace Web_Blazor.Services;

public class AuthState
{
    public bool IsAuthenticated { get; private set; }
    public string? Email { get; private set; }
    public string? Password { get; private set; }

    public event Action? OnChange;

    public bool TryLogin(string email, string password)
    {
        if (email == "test@test.com" && password == "0000")
        {
            Email = email;
            Password = password;
            IsAuthenticated = true;
            NotifyStateChanged();
            return true;
        }

        return false;
    }

    public void Logout()
    {
        Email = null;
        Password = null;
        IsAuthenticated = false;
        NotifyStateChanged();
    }

    private void NotifyStateChanged() => OnChange?.Invoke();
}
