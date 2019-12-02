import numpy as np
import matplotlib.pyplot as plt
import pandas as pd
import sys
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn.preprocessing import PolynomialFeatures

# cargar dataset
dataset = pd.read_csv('../machine_learning/strawberriesData.csv')

# seleccion de datos
X = dataset.iloc[:, 5].values.reshape(-1,1)
y = dataset.iloc[:,1].values

# Separacion de dataset (entrenamiento y prueba)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=0)

# Regresion Lineal
lin_reg = LinearRegression()
lin_reg.fit(X, y)

# Regresion Polinomial de grado 3
poly_reg = PolynomialFeatures(degree=3)
X_poly = poly_reg.fit_transform(X)
pol_reg = LinearRegression()
pol_reg.fit(X_poly, y)

def makePrediction(toPred):
    # Prediccion de modelo lineal y polinomial
    print(lin_reg.predict([[toPred]]), pol_reg.predict(poly_reg.fit_transform([[toPred]])), sep=",")

if __name__ == '__main__':
    globals()[sys.argv[1]](float(sys.argv[2]))
