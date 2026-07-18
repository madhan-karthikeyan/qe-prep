FROM maven:3.9-eclipse-temurin-21-alpine AS builder
WORKDIR /app
COPY pom.xml ./
COPY */pom.xml ./
RUN mvn dependency:go-offline -q
COPY . .
RUN mvn package -q -DskipTests

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /app/*/target/*.jar ./
CMD ["java", "--version"]
