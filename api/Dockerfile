FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS base
WORKDIR /app

FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src
COPY ["api.csproj", "./"]
RUN dotnet restore "api.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "api.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "api.csproj" -c Release -o /app/publish

FROM base AS final
RUN adduser --disabled-password --no-create-home --gecos "" service-user
USER service-user
WORKDIR /app
COPY --chown=service-user:service-user --from=publish /app/publish .
EXPOSE 8080
ENTRYPOINT ["dotnet", "api.dll"]

# docker build -t api .
# docker run -p 8080:8080 api