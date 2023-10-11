package com.example;

/**
 * Hello world!
 */
public final class App {
    private App() {
    }

    /**
     * Says hello to the world.
     * @param args The arguments of the program.
     */
    public static void main(String[] args) throws InterruptedException {
        for (int i=0; i!=-1;i++){ // Inverted condition to create infinte loop
            System.out.println("{\"level\":\"info\",\"time\":\"2021-08-16T18:26:46.621Z\",\"name\":\"bme280.data\",\"msg\":\"data\",\"temperature\":26.03,\"pressure\":999.39,\"humidity\":45.32,\"altitude\":116}");
            Thread.sleep(300000);
        }
    }
}
